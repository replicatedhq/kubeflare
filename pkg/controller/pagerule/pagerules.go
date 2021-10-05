package pagerule

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/internal"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
)

func (r *ReconcilePageRule) ReconcilePageRules(ctx context.Context, instance crdsv1alpha1.PageRule, zone *crdsv1alpha1.Zone, cf *cloudflare.API) error {
	logger.Debug("reconcilePageRules for zone")

	containsFinalizer := controllerutil.ContainsFinalizer(&instance, internal.DeleteCFResourceFinalizer)

	if instance.DeletionTimestamp != nil && !containsFinalizer {
		// object is being deleted and finalizer already executed. nothing more to do
		return nil
	}

	zoneID, err := cf.ZoneIDByName(zone.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	if instance.DeletionTimestamp != nil && containsFinalizer {
		// object is being deleted execute finalizer
		return r.finalize(ctx, &instance, zoneID, cf)
	}

	// TODO(user) status vs spec. would need to compute this if its not set or find a way to save the id to survive k8s object restores
	// if this page rule does not exist, create it, update the status and return
	if instance.Status.ID == "" {
		err = r.createPageRule(ctx, &instance, zoneID, cf)
		if err != nil {
			return err
		}
	}

	if instance.Status.ID == "" {
		// not managed by us
		return nil
	}

	if instance.DeletionTimestamp == nil && !containsFinalizer {
		patch := client.MergeFrom(instance.DeepCopy())
		controllerutil.AddFinalizer(&instance, internal.DeleteCFResourceFinalizer)

		if err = r.Client.Patch(ctx, &instance, patch); err != nil {
			return errors.Wrap(err, "could not add finalizer to page rule")
		}

		logger.Debug("added finalizer to page rule",
			zap.String("name", instance.Name),
			zap.String("finalizer", internal.DeleteCFResourceFinalizer),
			zap.String("zone", instance.Spec.Zone),
			zap.String("requestUrl", instance.Spec.Rule.RequestURL))
	}

	return nil
}

func (r *ReconcilePageRule) createPageRule(ctx context.Context, instance *crdsv1alpha1.PageRule, zoneID string, cf *cloudflare.API) error {
	rule := cloudflare.PageRule{}

	if instance.Spec.Rule.Enabled {
		rule.Status = "active"
	} else {
		rule.Status = "disabled"
	}

	if instance.Spec.Rule.Priority != nil {
		rule.Priority = *instance.Spec.Rule.Priority
	}

	rule.Targets = []cloudflare.PageRuleTarget{
		{
			Target: "url",
			Constraint: struct {
				Operator string `json:"operator"`
				Value    string `json:"value"`
			}{
				Operator: "matches",
				Value:    instance.Spec.Rule.RequestURL,
			},
		},
	}

	if instance.Spec.Rule.ForwardingURL != nil {
		rule.Actions = []cloudflare.PageRuleAction{
			{
				ID: "forwarding_url",
				Value: map[string]interface{}{
					"url":         instance.Spec.Rule.ForwardingURL.RedirectURL,
					"status_code": instance.Spec.Rule.ForwardingURL.StatusCode,
				},
			},
		}
	}
	if instance.Spec.Rule.AlwaysUseHTTPS != nil {
		rule.Actions = []cloudflare.PageRuleAction{
			{
				ID: "always_use_https",
			},
		}
	}

	created, err := cf.CreatePageRule(zoneID, rule)
	if err != nil {
		// There is no clean way (something that does not involve looking for strings) at the moment to determine that the error is as a result of a conflict/duplicate
		if strings.Contains(err.Error(), "Your zone already has an existing page rule with that URL") {
			// TODO(user) import existing records or ignore?
			// This is a duplicate, the page rule was probably not created by the operator. We ignore and stop retrying, or, should we?
			logger.Warn("page rule already exist, ignoring..", zap.String("zone", instance.Spec.Zone), zap.String("requestUrl", instance.Spec.Rule.RequestURL))
			return nil
		}

		return errors.Wrap(err, "failed to create page rule")
	}

	instance.Status.ID = created.ID
	if err := r.Status().Update(ctx, instance); err != nil {
		return errors.Wrap(err, "failed to update status")
	}

	logger.Debug("created page rule", zap.String("id", created.ID), zap.String("zone", instance.Spec.Zone), zap.String("requestUrl", instance.Spec.Rule.RequestURL))
	return nil
}

func (r *ReconcilePageRule) finalize(ctx context.Context, instance *crdsv1alpha1.PageRule, zoneID string, cf *cloudflare.API) error {
	if err := cf.DeletePageRule(zoneID, instance.Status.ID); err != nil {
		return errors.Wrap(err, "failed to delete page rule")
	}

	patch := client.MergeFrom(instance.DeepCopy())
	controllerutil.RemoveFinalizer(instance, internal.DeleteCFResourceFinalizer)

	if err := client.IgnoreNotFound(r.Client.Patch(ctx, instance, patch)); err != nil {
		return errors.Wrap(err, "failed to remove finalizer from page rule")
	}

	logger.Debug("removed page rule", zap.String("zone", instance.Spec.Zone), zap.String("id", instance.Status.ID), zap.String("requestUrl", instance.Spec.Rule.RequestURL))
	return nil
}
