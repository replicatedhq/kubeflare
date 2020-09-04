package shared

import (
	"context"

	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	crdsclientv1alpha1 "github.com/replicatedhq/kubeflare/pkg/client/kubeflareclientset/typed/crds/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func GetZone(ctx context.Context, namespace string, zoneName string) (*crdsv1alpha1.Zone, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config")
	}

	crdsClient, err := crdsclientv1alpha1.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create crds client")
	}

	zone, err := crdsClient.Zones(namespace).Get(ctx, zoneName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get zone")
	}

	return zone, nil
}
