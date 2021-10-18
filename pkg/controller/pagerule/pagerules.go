package pagerule

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/internal"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcilePageRule) ReconcilePageRules(ctx context.Context, instance crdsv1alpha1.PageRule, zone *crdsv1alpha1.Zone, cf *cloudflare.API) error {
	logger.Debug("reconcilePageRules for zone")

	containsFinalizer := controllerutil.ContainsFinalizer(&instance, internal.DeleteCFResourceFinalizer)
	isBeingDeleted := !instance.DeletionTimestamp.IsZero()

	if isBeingDeleted && !containsFinalizer {
		// object is being deleted and finalizer already executed. nothing more to do
		return nil
	}

	zoneID, err := cf.ZoneIDByName(zone.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	if isBeingDeleted && containsFinalizer {
		// object is being deleted execute finalizer
		return r.finalize(ctx, &instance, zoneID, cf)
	}

	justCreated := false

	// TODO(user) status vs spec. would need to compute this if its not set or find a way to save the id to survive k8s object restores
	// if this page rule does not exist, create it, update the status and return
	if instance.Status.ID == "" {
		err = r.createPageRule(ctx, &instance, zoneID, cf)

		if err != nil {
			// There is no clean way (something that does not involve looking for strings) at the moment to determine that the error is as a result of a conflict/duplicate
			if strings.Contains(err.Error(), "Your zone already has an existing page rule with that URL") {
				// TODO(user) import existing records or ignore?
				// This is a duplicate, the page rule was probably not created by the operator. We ignore and stop retrying, or, should we?
				logger.Warn("page rule already exist, ignoring..", zap.String("zone", instance.Spec.Zone), zap.String("requestUrl", instance.Spec.Rule.RequestURL))
				return nil
			}

			return err
		}

		justCreated = true
	}

	if !isBeingDeleted && !containsFinalizer {
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

	if justCreated {
		return nil
	}

	return r.syncPageRule(ctx, &instance, zoneID, cf)
}

func (r *ReconcilePageRule) createPageRule(ctx context.Context, instance *crdsv1alpha1.PageRule, zoneID string, cf *cloudflare.API) error {
	pageRule := r.mapCRDToCF(instance)
	created, err := cf.CreatePageRule(zoneID, pageRule)
	if err != nil {
		return errors.Wrap(err, "failed to create page rule")
	}

	created.Priority = pageRule.Priority

	instance.Status.ID = created.ID
	// We therefore store the user provided priority in the status and use that to check for differences later
	instance.Status.LastAppliedPriority = pageRule.Priority
	if err != nil {
		return errors.Wrap(err, "failed to compute page rule checksum")
	}

	if err := r.Status().Update(ctx, instance); err != nil {
		return errors.Wrap(err, "failed to update status")
	}

	logger.Debug("created page rule", zap.String("id", created.ID), zap.String("zone", instance.Spec.Zone), zap.String("requestUrl", instance.Spec.Rule.RequestURL))
	return nil
}

func (r *ReconcilePageRule) syncPageRule(ctx context.Context, instance *crdsv1alpha1.PageRule, zoneID string, cf *cloudflare.API) error {
	ruleCF, err := cf.PageRule(zoneID, instance.Status.ID)
	if err != nil {
		return errors.Wrap(err, "could not retrieve page rule")
	}

	ruleK8s := r.mapCRDToCF(instance)

	// Use the priority value set on the status as the "previous" value
	ruleCF.Priority = instance.Status.LastAppliedPriority

	if r.pageRulesAreEqual(ruleCF, ruleK8s) {
		return nil
	}

	if err = cf.UpdatePageRule(zoneID, ruleK8s.ID, ruleK8s); err != nil {
		return errors.Wrap(err, "could not update page rule")
	}

	// The priority might change because cloudflare removes gaps from priorities
	// We therefore store the user provided priority in the status and use that to check for differences later
	instance.Status.LastAppliedPriority = ruleK8s.Priority
	if err := r.Status().Update(ctx, instance); err != nil {
		return errors.Wrap(err, "failed to update status")
	}

	logger.Debug("updated page rule", zap.String("id", ruleK8s.ID), zap.String("zone", instance.Spec.Zone), zap.String("requestUrl", instance.Spec.Rule.RequestURL))
	return nil
}

func (r *ReconcilePageRule) deletePageRule(instance *crdsv1alpha1.PageRule, zoneID string, cf *cloudflare.API) error {
	err := cf.DeletePageRule(zoneID, instance.Status.ID)
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "Invalid Page Rule identifier") {
		// page rule no longer exist
		logger.Debug("page rule already deleted from cloudflare", zap.String("zone", instance.Spec.Zone), zap.String("id", instance.Status.ID), zap.String("requestUrl", instance.Spec.Rule.RequestURL))
		return nil
	}

	return err
}

func (r *ReconcilePageRule) finalize(ctx context.Context, instance *crdsv1alpha1.PageRule, zoneID string, cf *cloudflare.API) error {
	if err := r.deletePageRule(instance, zoneID, cf); err != nil {
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

func (r *ReconcilePageRule) mapCRDToCF(instance *crdsv1alpha1.PageRule) cloudflare.PageRule {
	rule := cloudflare.PageRule{ID: instance.Status.ID}

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

	if instance.Spec.Rule.CacheLevel != nil {
		rule.Actions = append(rule.Actions, cloudflare.PageRuleAction{
			ID:    "cache_level",
			Value: instance.Spec.Rule.CacheLevel.Level,
		})
	}

	if instance.Spec.Rule.EdgeCacheTTL != nil {
		rule.Actions = append(rule.Actions, cloudflare.PageRuleAction{
			ID:    "edge_cache_ttl",
			Value: instance.Spec.Rule.EdgeCacheTTL.Value,
		})
	}

	if instance.Spec.Rule.CacheKeyFields != nil {
		prepareArray := func(ar []string) []string {
			res := make([]string, len(ar))
			copy(res, ar)
			sort.Strings(res)
			return res
		}

		getQueryStringMap := func(field crdsv1alpha1.CacheKeyQueryStringField) map[string]interface{} {
			if field.Ignore {
				return map[string]interface{}{
					"include": []string{},
					"exclude": "*",
				}
			}

			if len(field.Include) == 0 && len(field.Exclude) == 0 {
				return map[string]interface{}{
					"include": "*",
					"exclude": []string{},
				}
			}

			if len(field.Include) == 1 && field.Include[0] == "*" {
				return map[string]interface{}{
					"include": "*",
					"exclude": []string{},
				}
			}

			if len(field.Exclude) == 1 && field.Exclude[0] == "*" {
				return map[string]interface{}{
					"include": []string{},
					"exclude": "*",
				}
			}

			return map[string]interface{}{
				"include": prepareArray(field.Include),
				"exclude": prepareArray(field.Exclude),
			}
		}

		rule.Actions = append(rule.Actions, cloudflare.PageRuleAction{
			ID: "cache_key_fields",
			Value: map[string]interface{}{
				"cookie": map[string][]string{
					"check_presence": prepareArray(instance.Spec.Rule.CacheKeyFields.Cookie.CheckPresence),
					"include":        prepareArray(instance.Spec.Rule.CacheKeyFields.Cookie.Include),
				},
				"header": map[string][]string{
					"check_presence": prepareArray(instance.Spec.Rule.CacheKeyFields.Header.CheckPresence),
					"include":        prepareArray(instance.Spec.Rule.CacheKeyFields.Header.Include),
					"exclude":        prepareArray(instance.Spec.Rule.CacheKeyFields.Header.Exclude),
				},
				"host": map[string]bool{
					"resolved": instance.Spec.Rule.CacheKeyFields.Host.Resolved,
				},
				"query_string": getQueryStringMap(instance.Spec.Rule.CacheKeyFields.QueryString),
				"user": map[string]bool{
					"device_type": instance.Spec.Rule.CacheKeyFields.User.DeviceType,
					"geo":         instance.Spec.Rule.CacheKeyFields.User.Geo,
					"lang":        instance.Spec.Rule.CacheKeyFields.User.Lang,
				},
			},
		})
	}

	if instance.Spec.Rule.ExplicitCacheControl != nil {
		rule.Actions = append(rule.Actions, cloudflare.PageRuleAction{
			ID:    "explicit_cache_control",
			Value: instance.Spec.Rule.ExplicitCacheControl.Value,
		})
	}

	if instance.Spec.Rule.SortQueryStrings != nil {
		rule.Actions = append(rule.Actions, cloudflare.PageRuleAction{
			ID:    "sort_query_string_for_cache",
			Value: instance.Spec.Rule.SortQueryStrings.Value,
		})
	}

	if instance.Spec.Rule.CacheOnCookie != nil {
		rule.Actions = append(rule.Actions, cloudflare.PageRuleAction{
			ID:    "cache_on_cookie",
			Value: instance.Spec.Rule.CacheOnCookie.Value,
		})
	}

	// overwrite everything else because cloudflare does not allow the forwarding_url action with any other action.
	// We'll add some validations for these later
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

	// overwrite everything else because cloudflare does not allow the always_use_https action with any other action.
	if instance.Spec.Rule.AlwaysUseHTTPS != nil {
		rule.Actions = []cloudflare.PageRuleAction{
			{
				ID: "always_use_https",
			},
		}
	}

	return rule
}

func (r *ReconcilePageRule) ruleActionValuesAreEqual(action1, action2 interface{}) bool {
	data1, err := json.Marshal(action1)
	if err != nil {
		logger.Error(err)
		return false
	}

	data2, err := json.Marshal(action2)
	if err != nil {
		logger.Error(err)
		return false
	}

	return string(data1) == string(data2)
}

func (r *ReconcilePageRule) pageRulesAreEqual(rule1, rule2 cloudflare.PageRule) bool {
	if rule1.Status != rule2.Status {
		logger.Debug("page rule status is not in sync", zap.String("status1", rule1.Status), zap.String("status2", rule2.Status))
		return false
	}

	if rule1.Priority != rule2.Priority {
		logger.Debug("page rule priority is not in sync", zap.Int("priority1", rule1.Priority), zap.Int("priority2", rule2.Priority))
		return false
	}

	numActions := len(rule1.Actions)
	numActions2 := len(rule2.Actions)
	if numActions != numActions2 {
		logger.Debug("page rule actions are not in sync", zap.Int("rule1NumberOfActions", numActions), zap.Int("rule1NumberOfActions", numActions2))
		return false
	}

	actionsMap1, actionsMap2 := make(map[string]interface{}, len(rule1.Actions)), make(map[string]interface{}, len(rule1.Actions))
	for i, action1 := range rule1.Actions {
		action2 := rule2.Actions[i]

		// IDs could be unequal due to actions array sort order
		if action2.ID != action1.ID {
			actionsMap1[action1.ID] = action1.Value
			actionsMap2[action2.ID] = action2.Value
			continue
		}

		// Tried reflect.DeepEqual(action1.Value, action2.Value) but it returns false
		// even when both interfaces represent the same concrete value. I think checking the underlying
		// type of the action value object, casting to a concrete type and then comparing would mean there is yet
		// another place to update when making changes to supported cloudflare page rule actions
		if !r.ruleActionValuesAreEqual(action1.Value, action2.Value) {
			// found a rule with matching id. check and fail fast if not equal. no need to add to map
			logger.Debug("page rule actions are not in sync", zap.String("actionID", action1.ID), zap.Any("action1", action1.Value), zap.Any("action2", action2.Value))
			return false
		}

		actionsMap1[action1.ID] = action1.Value
		actionsMap2[action2.ID] = action2.Value
	}

	for id, value1 := range actionsMap1 {
		value2, found := actionsMap2[id]

		if !found {
			logger.Debug("page rule actions are not in sync. action id not present in both rule sets", zap.String("actionID", id))
			return false
		}

		if !r.ruleActionValuesAreEqual(value1, value2) {
			logger.Debug("page rule actions are not in sync", zap.String("actionID", id), zap.Any("action1", value1), zap.Any("action2", value2))
			return false
		}
	}

	return true
}
