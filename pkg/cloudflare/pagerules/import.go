package pagerules

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	"github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FetchPageRulesForZone(token string, zone string, zoneID string) ([]*v1alpha1.PageRule, error) {
	cf, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "create clouflare client")
	}

	resources, err := cf.ListPageRules(zoneID)
	if err != nil {
		return nil, errors.Wrap(err, "fetch resources")
	}

	pageRules := []*v1alpha1.PageRule{}
	for _, resource := range resources {
		spec := v1alpha1.PageRuleSpec{
			Zone: zone,
		}

		for _, action := range resource.Actions {
			if action.ID == "forwarding_url" {
				spec.Rule = &v1alpha1.Rule{
					ForwardingURL: &v1alpha1.ForwardingURLPageRule{
						RedirectURL: action.Value.(map[string]interface{})["url"].(string),
						StatusCode:  int(action.Value.(map[string]interface{})["status_code"].(float64)),
					},
				}
			} else if action.ID == "always_use_https" {
				spec.Rule = &v1alpha1.Rule{
					AlwaysUseHTTPS: &v1alpha1.AlwaysUseHTTPSPageRule{},
				}
			}
		}

		if spec.Rule == nil {
			// we don't support this type of rule yet
			// prob should do something like add it
			continue
		}

		if resource.Status == "enabled" {
			spec.Rule.Enabled = true
		}

		spec.Rule.Priority = &resource.Priority

		for _, target := range resource.Targets {
			if target.Target == "url" {
				// this is a poor and partial implementation todo
				spec.Rule.RequestURL = target.Constraint.Value
			}
		}

		pageRule := v1alpha1.PageRule{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "crds.kubeflare.io/v1alpha1",
				Kind:       "PageRule",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: resource.ID, // TODO, this could be better
			},
			Spec: spec,
			Status: v1alpha1.PageRuleStatus{
				ID: resource.ID,
			},
		}

		pageRules = append(pageRules, &pageRule)
	}

	return pageRules, nil
}
