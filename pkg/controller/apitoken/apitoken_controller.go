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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
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

// Add creates a new APIToken Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	if err := add(mgr, newReconciler(mgr)); err != nil {
		return err
	}

	return addSecret(mgr, newSecretReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAPIToken{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("apitoken-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to APIToken
	err = c.Watch(&source.Kind{
		Type: &crdsv1alpha1.APIToken{},
	}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return errors.Wrap(err, "failed to start watch on apitokens")
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

var _ reconcile.Reconciler = &ReconcileAPIToken{}

// ReconcileAPIToken reconciles a APIToken object
type ReconcileAPIToken struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a APIToken object and makes changes based on the state read
// and what is in the APIToken.Spec
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=apitokens,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileAPIToken) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// This reconcile loop will be called for all APIToken objects
	// because of the informer that we have set up

	ctx := context.Background()
	instance := &crdsv1alpha1.APIToken{}

	if err := r.Get(ctx, request.NamespacedName, instance); err != nil {
		if apiErrors.IsNotFound(err) {
			logger.Debug("apitoken already deleted", zap.String("name", request.Name))
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

func (r *ReconcileAPIToken) handleDeletion(ctx context.Context, instance *crdsv1alpha1.APIToken) error {
	containsFinalizer := controllerutil.ContainsFinalizer(instance, internal.ProtectAPITokenFinalizer)
	isBeingDeleted := !instance.DeletionTimestamp.IsZero()

	if isBeingDeleted && !containsFinalizer {
		// object is being deleted and finalizer already executed. nothing more to do
		return nil
	}

	if isBeingDeleted && containsFinalizer {
		// object is being deleted check for dependencies
		hasDeps, err := r.apiTokenHasDependents(ctx, instance)
		if err != nil {
			return err
		}

		if hasDeps {
			logger.Debug("apitoken has dependencies and cannot be deleted yet", zap.String("name", instance.Name))
			return shared.HasDependenciesError
		}

		patch := client.MergeFrom(instance.DeepCopy())
		controllerutil.RemoveFinalizer(instance, internal.ProtectAPITokenFinalizer)

		if err := client.IgnoreNotFound(r.Client.Patch(ctx, instance, patch)); err != nil {
			return errors.Wrap(err, "failed to remove finalizer from apitoken")
		}

		logger.Debug("removed apitoken", zap.String("name", instance.Name))
		return nil
	}

	if !isBeingDeleted && !containsFinalizer {
		patch := client.MergeFrom(instance.DeepCopy())
		controllerutil.AddFinalizer(instance, internal.ProtectAPITokenFinalizer)

		if err := r.Client.Patch(ctx, instance, patch); err != nil {
			return errors.Wrap(err, "could not add finalizer to apitoken")
		}

		logger.Debug("added finalizer to apitoken",
			zap.String("name", instance.Name),
			zap.String("finalizer", internal.ProtectAPITokenFinalizer))

		if instance.Spec.ValueFrom != nil {
			// just in case the apitoken was created after the secret
			return r.ensureSecretHasFinalizer(ctx, instance.Spec.ValueFrom.SecretKeyRef.Name, instance.Namespace)
		}
	}

	return nil
}

func (r *ReconcileAPIToken) ensureSecretHasFinalizer(ctx context.Context, secretName, secretNamespace string) error {
	instance := &v1.Secret{}

	err := r.Get(ctx, types.NamespacedName{
		Namespace: secretNamespace,
		Name:      secretName,
	}, instance)

	if err != nil {
		if apiErrors.IsNotFound(err) {
			logger.Debug("referenced secret not found", zap.String("name", secretName), zap.String("namespace", secretNamespace))
			return nil
		}

		return err
	}

	if controllerutil.ContainsFinalizer(instance, internal.ProtectAPITokenFinalizer) {
		return nil
	}

	patch := client.MergeFrom(instance.DeepCopy())
	controllerutil.AddFinalizer(instance, internal.ProtectAPITokenFinalizer)

	if err := r.Client.Patch(ctx, instance, patch); err != nil {
		return errors.Wrap(err, "could not add finalizer to apitoken secret")
	}

	logger.Debug("added finalizer to apitoken secret",
		zap.String("name", secretName),
		zap.String("namespace", secretNamespace),
		zap.String("finalizer", internal.ProtectAPITokenFinalizer))

	return nil
}

func (r *ReconcileAPIToken) apiTokenHasDependents(ctx context.Context, apiToken *crdsv1alpha1.APIToken) (bool, error) {
	crdsClient, err := shared.GetCrdClient()
	if err != nil {
		return false, errors.Wrap(err, "failed to create crds client")
	}

	list, err := crdsClient.Zones(apiToken.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, zone := range list.Items {
		if zone.Spec.APIToken == apiToken.Name {
			return true, nil
		}
	}

	return false, nil
}
