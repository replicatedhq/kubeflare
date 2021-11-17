/*
Copyright 2021 The Kubernetes authors.

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

package workerroute

import (
	"context"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/controller/shared"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"go.uber.org/zap"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new WorkerRoute Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &WorkerRouteReconciler{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("workerroute-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to WorkerRoute
	err = c.Watch(&source.Kind{
		Type: &crdsv1alpha1.WorkerRoute{},
	}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return errors.Wrap(err, "failed to start watch on workerroutes")
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

// WorkerRouteReconciler reconciles a WorkerRoute object
type WorkerRouteReconciler struct {
	client.Client
	scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=workerroutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=workerroutes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=workerroutes/finalizers,verbs=update
func (r *WorkerRouteReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	var instance crdsv1alpha1.WorkerRoute
	if err := r.Get(ctx, request.NamespacedName, &instance); err != nil {
		if !apiErrors.IsNotFound(err) {
			logger.Error(err)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	zone, err := shared.GetZone(ctx, instance.Namespace, instance.Spec.Zone)
	if err != nil {
		logger.Error(err)
		return reconcile.Result{}, ignoreUnrecoverableErrors(err)
	}

	cf, err := shared.GetCloudflareAPI(ctx, instance.Namespace, zone.Spec.APIToken)
	if err != nil {
		logger.Error(err)
		return reconcile.Result{}, ignoreUnrecoverableErrors(err)
	}

	if err := r.ReconcileWorkerRouteInstances(ctx, &instance, zone, cf); err != nil {
		logger.Error(err)
		return reconcile.Result{}, ignoreUnrecoverableErrors(err)
	}

	return reconcile.Result{}, nil
}

func (r *WorkerRouteReconciler) ReconcileWorkerRouteInstances(ctx context.Context, instance *crdsv1alpha1.WorkerRoute, zone *crdsv1alpha1.Zone, cf *cloudflare.API) error {
	logger.Debug("reconcile worker route", zap.String("name", instance.Name))

	zoneID, err := cf.ZoneIDByName(zone.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	var finalizerKey = "crds.kubeflare.io/worker-route"
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !shared.ContainsString(instance.GetFinalizers(), finalizerKey) {
			controllerutil.AddFinalizer(instance, finalizerKey)
			if err := r.Update(ctx, instance); err != nil {
				return err
			}
		}
	} else {
		if shared.ContainsString(instance.GetFinalizers(), finalizerKey) {
			if err := r.deleteWorkerRoute(ctx, instance, zoneID, cf); err != nil {
				return err
			}
			controllerutil.RemoveFinalizer(instance, finalizerKey)
			if err := r.Update(ctx, instance); err != nil {
				return err
			}
		}
		return nil
	}

	if instance.Status.ID == "" {
		return r.createWorkerRoute(ctx, instance, zoneID, cf)
	}

	return r.updateWorkerRoute(ctx, instance, zoneID, cf)
}

func (r *WorkerRouteReconciler) createWorkerRoute(ctx context.Context, instance *crdsv1alpha1.WorkerRoute, zoneID string, cf *cloudflare.API) error {
	route := mapCRDToCF(instance)

	created, err := cf.CreateWorkerRoute(zoneID, route)
	if err != nil {
		if err := r.updateStatusLastError(ctx, instance, err); err != nil {
			return errors.Wrap(err, "failed to update status")
		}
		return errors.Wrap(err, "failed to create worker route")
	}
	logger.Debug("created worker route", zap.String("id", created.ID), zap.String("zone", instance.Spec.Zone), zap.String("pattern", instance.Spec.Pattern))

	instance.Status.ID = created.ID
	instance.Status.LastError = ""
	if err := r.Status().Update(ctx, instance); err != nil {
		return errors.Wrap(err, "failed to update status")
	}

	return nil
}

func (r *WorkerRouteReconciler) updateWorkerRoute(ctx context.Context, instance *crdsv1alpha1.WorkerRoute, zoneID string, cf *cloudflare.API) error {
	routeCF, err := getWorkerRoute(cf, zoneID, instance.Status.ID)
	if err != nil {
		if err := r.updateStatusLastError(ctx, instance, err); err != nil {
			return errors.Wrap(err, "failed to update status")
		}
		return err
	}

	route := mapCRDToCF(instance)
	if equals(routeCF, route) {
		logger.Debug("worker route is in sync", zap.String("name", instance.Name))
		return nil
	}

	updated, err := cf.UpdateWorkerRoute(zoneID, instance.Status.ID, route)
	if err != nil {
		if err := r.updateStatusLastError(ctx, instance, err); err != nil {
			return errors.Wrap(err, "failed to update status")
		}
		return errors.Wrap(err, "failed to update worker route")
	}
	logger.Debug("updated worker route", zap.String("id", updated.ID), zap.String("zone", instance.Spec.Zone), zap.String("pattern", instance.Spec.Pattern))

	instance.Status.LastError = ""
	if err := r.Status().Update(ctx, instance); err != nil {
		return errors.Wrap(err, "failed to update status")
	}
	return nil
}

func (r *WorkerRouteReconciler) deleteWorkerRoute(ctx context.Context, instance *crdsv1alpha1.WorkerRoute, zoneID string, cf *cloudflare.API) error {
	if instance.Status.ID == "" {
		return nil
	}
	deleted, err := cf.DeleteWorkerRoute(zoneID, instance.Status.ID)
	if err != nil {
		if err := r.updateStatusLastError(ctx, instance, err); err != nil {
			return errors.Wrap(err, "failed to update status")
		}
		return errors.Wrap(err, "failed to delete worker route")
	}
	logger.Debug("deleted worker route", zap.String("id", deleted.ID), zap.String("zone", instance.Spec.Zone), zap.String("pattern", instance.Spec.Pattern))
	return nil
}

func (r *WorkerRouteReconciler) updateStatusLastError(ctx context.Context, instance *crdsv1alpha1.WorkerRoute, err error) error {
	instance.Status.LastError = errors.Cause(err).Error()
	return r.Status().Update(ctx, instance)
}

func mapCRDToCF(instance *crdsv1alpha1.WorkerRoute) cloudflare.WorkerRoute {
	route := cloudflare.WorkerRoute{
		ID:      instance.Status.ID,
		Pattern: instance.Spec.Pattern,
		Script:  instance.Spec.Script,
	}
	return route
}

func ignoreUnrecoverableErrors(err error) error {
	if strings.Contains(err.Error(), "workers.api.error.duplicate_route") {
		return nil
	}
	if strings.Contains(err.Error(), "workers.api.error.invalid_route_script_missing") {
		return nil
	}
	if strings.Contains(err.Error(), "Route pattern must include zone name") {
		return nil
	}
	return err
}

// TODO: replace with cloudflare.GetWorkerRoute() when it is implemented in the upstream
func getWorkerRoute(cf *cloudflare.API, zoneID, ID string) (cloudflare.WorkerRoute, error) {
	workerRoutes, err := cf.ListWorkerRoutes(zoneID)
	if err != nil {
		return cloudflare.WorkerRoute{}, errors.Wrap(err, "failed to list worker routes")
	}
	for _, wr := range workerRoutes.Routes {
		if wr.ID == ID {
			return wr, nil
		}
	}
	return cloudflare.WorkerRoute{}, errors.Errorf("worker route not found: %s", ID)
}

func equals(w1, w2 cloudflare.WorkerRoute) bool {
	return w1.Pattern == w2.Pattern && w1.Script == w2.Script
}
