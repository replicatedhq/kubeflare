package dnsrecord

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/logger"
)

func ReconcileDNSRecords(ctx context.Context, instance crdsv1alpha1.DNSRecord, zone *crdsv1alpha1.Zone, cf *cloudflare.API) error {
	logger.Debug("reconcileDNSRecords for zone")

	zoneID, err := cf.ZoneIDByName(zone.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	existingRecords, err := cf.DNSRecords(zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return errors.Wrap(err, "failed to list dns records")
	}

	desiredRecords := []*crdsv1alpha1.Record{}
	if instance.Spec.Record != nil {
		desiredRecords = append(desiredRecords, instance.Spec.Record)
	}
	if instance.Spec.Records != nil {
		desiredRecords = append(desiredRecords, instance.Spec.Records...)
	}

	recordsToCreate := []cloudflare.DNSRecord{}
	recordsToUpdate := []cloudflare.DNSRecord{}
	// recordsToDelete := []cloudflare.DNSRecord{}

	for _, existingRecord := range existingRecords {
		found := false
		for _, desiredRecord := range desiredRecords {
			if desiredRecord.Name == existingRecord.Name && desiredRecord.Type == existingRecord.Type {
				found = true
				isChanged := false

				if desiredRecord.Content != existingRecord.Content {
					isChanged = true
					existingRecord.Content = desiredRecord.Content
				}
				desiredTTL := 1
				if desiredRecord.TTL != nil {
					desiredTTL = *desiredRecord.TTL
				}
				if desiredTTL != existingRecord.TTL {
					isChanged = true
					existingRecord.TTL = desiredTTL
				}
				if desiredRecord.Priority != nil {
					if *desiredRecord.Priority != existingRecord.Priority {
						isChanged = true
						existingRecord.Priority = *desiredRecord.Priority
					}
				}
				if desiredRecord.Proxied != nil {
					if *desiredRecord.Proxied != existingRecord.Proxied {
						isChanged = true
						existingRecord.Proxied = *desiredRecord.Proxied
					}
				}

				if isChanged {
					recordsToUpdate = append(recordsToUpdate, existingRecord)
				}
			}
		}
		if !found {
			// TODO this feels dangerous, how can we opt-in to delete somehow to avoid erasing all records
			// recordsToDelete = append(recordsToDelete, existingRecord)
		}
	}

	for _, desiredRecord := range desiredRecords {
		found := false
		for _, existingRecord := range existingRecords {
			if existingRecord.Type == desiredRecord.Type && existingRecord.Name == desiredRecord.Name {
				found = true
				goto Found
			}
		}
	Found:
		if !found {
			recordToCreate := cloudflare.DNSRecord{
				Type:    desiredRecord.Type,
				Name:    desiredRecord.Name,
				Content: desiredRecord.Content,
			}
			if desiredRecord.TTL != nil {
				recordToCreate.TTL = *desiredRecord.TTL
			} else {
				recordToCreate.TTL = 1
			}

			if desiredRecord.Priority != nil {
				recordToCreate.Priority = *desiredRecord.Priority
			}
			if desiredRecord.Proxied != nil {
				recordToCreate.Proxied = *desiredRecord.Proxied
			}
			recordsToCreate = append(recordsToCreate, recordToCreate)
		}
	}

	for _, recordToCreate := range recordsToCreate {
		response, err := cf.CreateDNSRecord(zoneID, recordToCreate)
		if err != nil {
			return errors.Wrap(err, "failed to create dns record")
		}

		if !response.Success {
			return errors.New("non success when creating dns record")
		}
	}

	for _, recordToUpdate := range recordsToUpdate {
		rr := cloudflare.DNSRecord{
			Type:    recordToUpdate.Type,
			Name:    recordToUpdate.Name,
			Content: recordToUpdate.Content,
			TTL:     recordToUpdate.TTL,
			Proxied: recordToUpdate.Proxied,
		}

		err := cf.UpdateDNSRecord(zoneID, recordToUpdate.ID, rr)
		if err != nil {
			return errors.Wrap(err, "failed to update dns record")
		}
	}

	return nil
}
