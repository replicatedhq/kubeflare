package dns

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	"github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FetchDNSRecordsForZone(token string, zone string, zoneID string) ([]*v1alpha1.DNSRecord, error) {
	cf, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "create clouflare client")
	}

	resources, err := cf.DNSRecords(zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return nil, errors.Wrap(err, "fetch resources")
	}

	dnsRecords := []*v1alpha1.DNSRecord{}
	for _, resource := range resources {
		dnsRecord := v1alpha1.DNSRecord{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "crds.kubeflare.io/v1alpha1",
				Kind:       "DNSRecord",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: resource.Name,
			},
			Spec: v1alpha1.DNSRecordSpec{
				Zone: zone,
				Record: &v1alpha1.Record{
					Type:     resource.Type,
					Name:     resource.Name,
					Content:  resource.Content,
					TTL:      &resource.TTL,
					Priority: &resource.Priority,
					Proxied:  &resource.Proxied,
				},
			},
		}

		dnsRecords = append(dnsRecords, &dnsRecord)
	}

	return dnsRecords, nil
}
