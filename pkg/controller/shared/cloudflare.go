package shared

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsclientv1alpha1 "github.com/replicatedhq/kubeflare/pkg/client/kubeflareclientset/typed/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var HasDependenciesError = errors.New("dependency detected")

func GetCloudflareAPI(ctx context.Context, namespace string, apiTokenName string) (*cloudflare.API, error) {
	crdsClient, err := GetCrdClient()
	if err != nil {
		return nil, err
	}

	apiToken, err := crdsClient.APITokens(namespace).Get(ctx, apiTokenName, metav1.GetOptions{})
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

	api, err := cloudflare.NewWithAPIToken(tokenValue)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloudflare api instance")
	}

	return api, nil
}

func GetCrdClient() (*crdsclientv1alpha1.CrdsV1alpha1Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config")
	}

	crdsClient, err := crdsclientv1alpha1.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create crds client")
	}

	return crdsClient, nil
}
