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
	"time"

	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/controller/shared"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"github.com/spf13/viper"
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
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	v := viper.GetViper()

	pollInterval, _ := time.ParseDuration(v.GetString("poll-interval"))

	return &ReconcileZone{
		Client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		pollInterval: pollInterval,
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
	scheme       *runtime.Scheme
	pollInterval time.Duration
}

// Reconcile reads that state of the cluster for a Zone object and makes changes based on the state read
// and what is in the Zone.Spec
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=zones,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=zones/status,verbs=get;update;patch
func (r *ReconcileZone) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// This reconcile loop will be called for all Zone objects
	// because of the informer that we have set up
	ctx := context.Background()
	instance := crdsv1alpha1.Zone{}
	err := r.Get(ctx, request.NamespacedName, &instance)
	if err != nil {
		logger.Error(err)
		return reconcile.Result{}, err
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

	return reconcile.Result{RequeueAfter: r.pollInterval}, nil
}
