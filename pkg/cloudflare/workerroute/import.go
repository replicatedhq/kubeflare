package workerroute

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	"github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FetchWorkerRoutesForZone(token string, zone string, zoneID string) ([]*v1alpha1.WorkerRoute, error) {
	cf, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "create clouflare client")
	}

	resources, err := cf.ListWorkerRoutes(context.Background(), zoneID)
	if err != nil {
		return nil, errors.Wrap(err, "fetch resources")
	}

	workerRoutes := []*v1alpha1.WorkerRoute{}
	for i, resource := range resources.Routes {
		workerRoute := v1alpha1.WorkerRoute{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "crds.kubeflare.io/v1alpha1",
				Kind:       "WorkerRoute",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("workerroute-%d", i),
				Annotations: map[string]string{
					internal.ImportedIDAnnotation: resource.ID,
				},
			},
			Spec: v1alpha1.WorkerRouteSpec{
				Zone:    zone,
				Pattern: resource.Pattern,
				Script:  resource.Script,
			},
		}

		workerRoutes = append(workerRoutes, &workerRoute)
	}

	return workerRoutes, nil
}
