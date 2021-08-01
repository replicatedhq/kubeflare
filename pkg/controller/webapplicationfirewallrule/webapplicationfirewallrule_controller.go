/*
Copyright 2020 Replicated, Inc.

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

package webapplicationfirewallrule

import (
	"context"
	"time"

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
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("webapplicationfirewallrule-controller")

// Add creates a new WebApplicationFirewallRule Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	v := viper.GetViper()

	pollInterval, _ := time.ParseDuration(v.GetString("poll-interval"))

	return &ReconcileWebApplicationFirewallRule{
		Client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		pollInterval: pollInterval,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("webapplicationfirewallrule-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to WebApplicationFirewallRule
	err = c.Watch(&source.Kind{Type: &crdsv1alpha1.WebApplicationFirewallRule{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
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

var _ reconcile.Reconciler = &ReconcileWebApplicationFirewallRule{}

// ReconcileWebApplicationFirewallRule reconciles a WebApplicationFirewallRule object
type ReconcileWebApplicationFirewallRule struct {
	client.Client
	scheme       *runtime.Scheme
	pollInterval time.Duration
}

// Reconcile reads that state of the cluster for a WebApplicationFirewallRule object and makes changes based on the state read
// and what is in the WebApplicationFirewallRule.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=webapplicationfirewallrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crds.kubeflare.io,resources=webapplicationfirewallrules/status,verbs=get;update;patch
func (r *ReconcileWebApplicationFirewallRule) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the WebApplicationFirewallRule instance
	ctx := context.Background()
	instance := crdsv1alpha1.WebApplicationFirewallRule{}
	err := r.Get(ctx, request.NamespacedName, &instance)
	if err != nil {
		logger.Error(err)
		return reconcile.Result{}, err
	}

	zone, err := shared.GetZone(ctx, instance.Namespace, instance.Spec.Zone)
	if err != nil {
		logger.Error(err)
		return reconcile.Result{}, err
	}

	cf, err := shared.GetCloudflareAPI(ctx, instance.Namespace, zone.Spec.APIToken)
	if err != nil {
		logger.Error(err)
		return reconcile.Result{}, err
	}

	if err := ReconcileWAFRuleInstances(ctx, instance, zone, cf); err != nil {
		logger.Error(err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: r.pollInterval}, nil
}
