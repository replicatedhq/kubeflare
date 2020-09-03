package zone

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	crdsclientv1alpha1 "github.com/replicatedhq/kubeflare/pkg/client/kubeflareclientset/typed/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func (r *ReconcileZone) getCloudflareAPI(ctx context.Context, instance crdsv1alpha1.Zone) (*cloudflare.API, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config")
	}

	crdsClient, err := crdsclientv1alpha1.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create crds client")
	}

	apiToken, err := crdsClient.APITokens(instance.Namespace).Get(ctx, instance.Spec.APIToken, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get api token")
	}

	tokenValue, err := apiToken.GetTokenValue(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token value")
	}

	logger.Debug("creating cloudflare api object",
		zap.String("email", apiToken.Spec.Email),
		zap.Int("tokenLength", len(tokenValue)))

	api, err := cloudflare.New(tokenValue, apiToken.Spec.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloudflare api instance")
	}

	return api, nil
}
