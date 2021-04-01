package webapplicationfirewallrule

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/logger"
)

func ReconcileWAFRuleInstances(ctx context.Context, instance crdsv1alpha1.WebApplicationFirewallRule, zone *crdsv1alpha1.Zone, cf *cloudflare.API) error {
	logger.Debug("ReconcileWAFRules for zone")

	zoneID, err := cf.ZoneIDByName(zone.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	existingPackages, err := cf.ListWAFPackages(zoneID)
	if err != nil {
		return errors.Wrap(err, "failed to list WAF packages")
	}

	existingRules := []cloudflare.WAFRule{}
	for _, currentPackage := range existingPackages {
		rules, err := cf.ListWAFRules(zoneID, currentPackage.ID)
		if err != nil {
			return errors.Wrap(err, "failed to list WAF rules")
		}

		existingRules = append(existingRules, rules...)
	}

	desiredRules := []*crdsv1alpha1.WAFRule{}
	if instance.Spec.Rule != nil {
		desiredRules = append(desiredRules, instance.Spec.Rule)
	}
	if instance.Spec.Rules != nil {
		desiredRules = append(desiredRules, instance.Spec.Rules...)
	}

	rulesToUpdate := []cloudflare.WAFRule{}

	for _, existingRule := range existingRules {
		found := false
		for _, desiredRule := range desiredRules {
			if desiredRule.ID == existingRule.ID {
				found = true
				isChanged := false

				if desiredRule.Mode != existingRule.Mode {
					isChanged = true
					existingRule.Mode = desiredRule.Mode
				}

				if isChanged {
					rulesToUpdate = append(rulesToUpdate, existingRule)
				}
			}
		}
		if !found {
			// TODO
		}
	}

	for _, ruleToUpdate := range rulesToUpdate {
		_, err := cf.UpdateWAFRule(zoneID, ruleToUpdate.PackageID, ruleToUpdate.ID, ruleToUpdate.Mode)
		if err != nil {
			return errors.Wrap(err, "failed to update WAF rule")
		}
	}

	return nil
}
