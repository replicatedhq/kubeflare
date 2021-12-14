/*
Copyright 2019 Replicated, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apitoken

import (
	"context"
	"github.com/replicatedhq/kubeflare/pkg/internal"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedhq/kubeflare/pkg/controller/shared"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"k8s.io/apimachinery/pkg/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// newSecretReconciler returns a new reconcile.Reconciler
func newSecretReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &reconcileSecret{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
}

// addSecret adds a new Controller to mgr with r as the reconcile.Reconciler
func addSecret(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("secret-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Secret
	err = c.Watch(&source.Kind{
		Type: &v1.Secret{},
	}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return errors.Wrap(err, "failed to start watch on secrets")
	}

	generatedClient := kubernetes.NewForConfigOrDie(mgr.GetConfig())
	generatedInformers := kubeinformers.NewSharedInformerFactory(generatedClient, time.Minute)
	err = mgr.Add(manager.RunnableFunc(func(s <-chan struct{}) error {
		generatedInformers.Start(s)
		<-s
		return nil
	}))
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &reconcileSecret{}

// ReconcileSecret reconciles a Secrets related to an APIToken object
type reconcileSecret struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Secret object and makes changes based on the state read
// if the secrets is referenced from an APIToken and protectAPIToken is true
// +kubebuilder:rbac:groups=,resources=secrets,verbs=create,update,patch,delete
func (r *reconcileSecret) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	instance := &v1.Secret{}

	if err := r.Get(ctx, request.NamespacedName, instance); err != nil {
		if apiErrors.IsNotFound(err) {
			logger.Debug("secret already deleted",
				zap.String("name", request.Name), zap.String("namespace", request.Namespace))

			return reconcile.Result{}, nil
		}

		logger.Error(err)
		return reconcile.Result{}, err
	}

	if err := r.handleDeletion(ctx, instance); err != nil {
		if errors.Is(err, shared.HasDependenciesError) {
			return reconcile.Result{
				RequeueAfter: time.Duration(10) * time.Second,
			}, nil
		}

		logger.Error(err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *reconcileSecret) handleDeletion(ctx context.Context, instance *v1.Secret) error {
	containsFinalizer := controllerutil.ContainsFinalizer(instance, internal.ProtectAPITokenFinalizer)
	isBeingDeleted := !instance.DeletionTimestamp.IsZero()

	if isBeingDeleted && !containsFinalizer {
		// object is being deleted and finalizer already executed. nothing more to do
		return nil
	}

	if isBeingDeleted && containsFinalizer {
		// object is being deleted check for dependencies
		hasDeps, err := secretHasDependents(ctx, instance)
		if err != nil {
			logger.Error(err)
			return err
		}

		if hasDeps {
			logger.Debug("secret has dependencies and cannot be deleted yet", zap.String("name", instance.Name))
			return shared.HasDependenciesError
		}

		patch := client.MergeFrom(instance.DeepCopy())
		controllerutil.RemoveFinalizer(instance, internal.ProtectAPITokenFinalizer)

		if err := client.IgnoreNotFound(r.Client.Patch(ctx, instance, patch)); err != nil {
			return errors.Wrap(err, "failed to remove finalizer from secret")
		}

		logger.Debug("removed finalizer from secret",
			zap.String("name", instance.Name), zap.String("namespace", instance.Namespace))

		return nil
	}

	if !isBeingDeleted && !containsFinalizer {
		hasDeps, err := secretHasDependents(ctx, instance)
		if err != nil {
			logger.Error(err)
			return err
		}

		if !hasDeps {
			logger.Debug("not adding finalizer to secret, not referenced by any apitoken",
				zap.String("name", instance.Name),
				zap.String("namespace", instance.Namespace),
				zap.String("finalizer", internal.ProtectAPITokenFinalizer))

			return nil
		}

		patch := client.MergeFrom(instance.DeepCopy())
		controllerutil.AddFinalizer(instance, internal.ProtectAPITokenFinalizer)

		if err := r.Client.Patch(ctx, instance, patch); err != nil {
			return errors.Wrap(err, "could not add finalizer to secret")
		}

		logger.Debug("added finalizer to secret",
			zap.String("name", instance.Name),
			zap.String("namespace", instance.Namespace),
			zap.String("finalizer", internal.ProtectAPITokenFinalizer))
	}

	return nil
}

func secretHasDependents(ctx context.Context, secret *v1.Secret) (bool, error) {
	crdsClient, err := shared.GetCrdClient()
	if err != nil {
		return false, err
	}

	list, err := crdsClient.APITokens(secret.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, apiToken := range list.Items {
		if apiToken.Spec.ValueFrom.SecretKeyRef.Name == secret.Name {
			return true, nil
		}
	}

	return false, nil
}
