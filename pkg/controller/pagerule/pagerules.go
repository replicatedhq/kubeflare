package pagerule

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/logger"
)

func ReconcilePageRules(ctx context.Context, instance crdsv1alpha1.PageRule, zone *crdsv1alpha1.Zone, cf *cloudflare.API) error {
	logger.Debug("reconcileDNSRecords for zone")

	return nil
}
