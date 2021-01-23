package pagerule

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/logger"
)

func (r *ReconcilePageRule) ReconcilePageRules(ctx context.Context, instance crdsv1alpha1.PageRule, zone *crdsv1alpha1.Zone, cf *cloudflare.API) error {
	logger.Debug("reconcilePageRules for zone")

	zoneID, err := cf.ZoneIDByName(zone.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	// if this page rule does not exist, create it, update the status and return
	if instance.Status.ID == "" {
		return r.createPageRule(ctx, instance, zoneID, cf)
	}

	return nil
}

func (r *ReconcilePageRule) createPageRule(ctx context.Context, instance crdsv1alpha1.PageRule, zoneID string, cf *cloudflare.API) error {
	rule := cloudflare.PageRule{}

	if instance.Spec.Rule.Enabled {
		rule.Status = "enabled"
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
		return errors.Wrap(err, "failed to create page rule")
	}

	instance.Status.ID = created.ID
	if err := r.Status().Update(ctx, &instance); err != nil {
		return errors.Wrap(err, "failed to update status")
	}

	return nil
}
