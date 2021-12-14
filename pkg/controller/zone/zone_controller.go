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

package zone

import (
	"context"
	crdsclientv1alpha1 "github.com/replicatedhq/kubeflare/pkg/client/kubeflareclientset/typed/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/internal"
	"go.uber.org/zap"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Add creates a new Zone Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, protectAPIToken bool) error {
	return add(mgr, newReconciler(mgr, protectAPIToken))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, protectAPIToken bool) reconcile.Reconciler {
	return &ReconcileZone{
		Client:          mgr.GetClient(),
		scheme:          mgr.GetScheme(),
		protectAPIToken: protectAPIToken,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("zone-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Zone
	err = c.Watch(&source.Kind{
		Type: &crdsv1alpha1.Zone{},
	}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return errors.Wrap(err, "failed to start watch on zones")
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

var _ reconcile.Reconciler = &ReconcileZone{}

// ReconcileZone reconciles a Zone object
type ReconcileZone struct {
	client.Client
	scheme          *runtime.Scheme
	protectAPIToken bool
}

// Reconcile reads that state of the cluster for a Zone object and makes changes based on the state read
// and what is in the Zone.Spec
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=zones,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=zones/status,verbs=get;update;patch
func (r *ReconcileZone) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// This reconcile loop will be called for all Zone objects
	// because of the informer that we have set up
	ctx := context.Background()
	instance := &crdsv1alpha1.Zone{}
	err := r.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			logger.Debug("zone already deleted", zap.String("name", request.Name))
			return reconcile.Result{}, nil
		}

		logger.Error(err)
		return reconcile.Result{}, err
	}

	if r.protectAPIToken {
		deleted, err := r.handleDeletion(ctx, instance)

		if errors.Is(err, shared.HasDependenciesError) {
			return reconcile.Result{
				RequeueAfter: time.Duration(10) * time.Second,
			}, nil
		}

		if err != nil {
			logger.Error(err)
			return reconcile.Result{}, err
		}

		if deleted {
			return reconcile.Result{}, nil
		}
	}

	cf, err := shared.GetCloudflareAPI(ctx, instance.Namespace, instance.Spec.APIToken)
	if err != nil {
		logger.Error(err)
		return reconcile.Result{}, err
	}

	if err := ReconcileSettings(ctx, instance, cf); err != nil {
		logger.Error(err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileZone) handleDeletion(ctx context.Context, instance *crdsv1alpha1.Zone) (bool, error) {
	containsFinalizer := controllerutil.ContainsFinalizer(instance, internal.ProtectAPITokenFinalizer)
	isBeingDeleted := !instance.DeletionTimestamp.IsZero()

	if isBeingDeleted && !containsFinalizer {
		// object is being deleted and finalizer already executed. nothing more to do
		return true, nil
	}

	if isBeingDeleted && containsFinalizer {
		// object is being deleted check for dependencies
		hasDeps, err := r.zoneHasDependents(ctx, instance)
		if err != nil {
			logger.Error(err)
			return false, err
		}

		if hasDeps {
			logger.Debug("zone has dependencies and cannot be deleted yet", zap.String("name", instance.Name))
			return false, shared.HasDependenciesError
		}

		patch := client.MergeFrom(instance.DeepCopy())
		controllerutil.RemoveFinalizer(instance, internal.ProtectAPITokenFinalizer)

		if err := client.IgnoreNotFound(r.Client.Patch(ctx, instance, patch)); err != nil {
			return false, errors.Wrap(err, "failed to remove finalizer from zone")
		}

		logger.Debug("removed zone", zap.String("name", instance.Name))
		return true, nil
	}

	if !isBeingDeleted && !containsFinalizer {
		patch := client.MergeFrom(instance.DeepCopy())
		controllerutil.AddFinalizer(instance, internal.ProtectAPITokenFinalizer)

		if err := r.Client.Patch(ctx, instance, patch); err != nil {
			return false, errors.Wrap(err, "could not add finalizer to zone")
		}

		logger.Debug("added finalizer to zone",
			zap.String("name", instance.Name),
			zap.String("finalizer", internal.ProtectAPITokenFinalizer))
	}

	return false, nil
}

func (r *ReconcileZone) zoneHasDependents(ctx context.Context, zone *crdsv1alpha1.Zone) (bool, error) {
	crdsClient, err := shared.GetCrdClient()
	if err != nil {
		return false, err
	}

	hasDeps, err := r.zoneHasAccessAppDependents(ctx, crdsClient.AccessApplications(zone.Namespace), zone.Name)
	if err != nil || hasDeps {
		return hasDeps, err
	}

	hasDeps, err = r.zoneHasDNSDependents(ctx, crdsClient.DNSRecords(zone.Namespace), zone.Name)
	if err != nil || hasDeps {
		return hasDeps, err
	}

	hasDeps, err = r.zoneHasPageRuleDependents(ctx, crdsClient.PageRules(zone.Namespace), zone.Name)
	if err != nil || hasDeps {
		return hasDeps, err
	}

	hasDeps, err = r.zoneHasWAFDependents(ctx, crdsClient.WebApplicationFirewallRules(zone.Namespace), zone.Name)
	if err != nil || hasDeps {
		return hasDeps, err
	}

	hasDeps, err = r.zoneHasWorkerRouteDependents(ctx, crdsClient.WorkerRoutes(zone.Namespace), zone.Name)
	if err != nil || hasDeps {
		return hasDeps, err
	}

	return false, nil
}

func (r *ReconcileZone) zoneHasAccessAppDependents(ctx context.Context, client crdsclientv1alpha1.AccessApplicationInterface, zone string) (bool, error) {
	list, err := client.List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, item := range list.Items {
		if item.Spec.Zone == zone {
			return true, nil
		}
	}

	return false, nil
}

func (r *ReconcileZone) zoneHasDNSDependents(ctx context.Context, client crdsclientv1alpha1.DNSRecordInterface, zone string) (bool, error) {
	list, err := client.List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, item := range list.Items {
		if item.Spec.Zone == zone {
			return true, nil
		}
	}

	return false, nil
}

func (r *ReconcileZone) zoneHasPageRuleDependents(ctx context.Context, client crdsclientv1alpha1.PageRuleInterface, zone string) (bool, error) {
	list, err := client.List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, item := range list.Items {
		if item.Spec.Zone == zone {
			return true, nil
		}
	}

	return false, nil
}

func (r *ReconcileZone) zoneHasWAFDependents(ctx context.Context, client crdsclientv1alpha1.WebApplicationFirewallRuleInterface, zone string) (bool, error) {
	list, err := client.List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, item := range list.Items {
		if item.Spec.Zone == zone {
			return true, nil
		}
	}

	return false, nil
}

func (r *ReconcileZone) zoneHasWorkerRouteDependents(ctx context.Context, client crdsclientv1alpha1.WorkerRouteInterface, zone string) (bool, error) {
	list, err := client.List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, item := range list.Items {
		if item.Spec.Zone == zone {
			return true, nil
		}
	}

	return false, nil
}
