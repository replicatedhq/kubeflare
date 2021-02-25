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

		// TODO parse actions

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
