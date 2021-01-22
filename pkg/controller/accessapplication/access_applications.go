package accessapplication

import (
	"context"
	"sort"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/logger"
)

func ReconcileAccessApplicationInstance(ctx context.Context, instance crdsv1alpha1.AccessApplication, zone *crdsv1alpha1.Zone, cf *cloudflare.API) (*cloudflare.AccessApplication, error) {
	logger.Debug("reconcileAccessApplication for zone")

	if instance.Status.ApplicationID == "" {
		accessApplication := cloudflare.AccessApplication{
			Name:   instance.Spec.Name,
			Domain: instance.Spec.Domain,
		}

		if instance.Spec.SessionDuration != "" {
			accessApplication.SessionDuration = instance.Spec.SessionDuration
		}
		if instance.Spec.AllowedIdPs != nil && len(instance.Spec.AllowedIdPs) > 0 {
			accessApplication.AllowedIdps = instance.Spec.AllowedIdPs
		}
		if instance.Spec.AutoRedirectToIdentity != nil {
			accessApplication.AutoRedirectToIdentity = *instance.Spec.AutoRedirectToIdentity
		}
		if instance.Spec.CORSHeaders != nil {
			corsHeaders := cloudflare.AccessApplicationCorsHeaders{
				AllowedMethods:   instance.Spec.CORSHeaders.AllowedMethods,
				AllowedOrigins:   instance.Spec.CORSHeaders.AllowedOrigins,
				AllowedHeaders:   instance.Spec.CORSHeaders.AllowedHeaders,
				AllowAllMethods:  instance.Spec.CORSHeaders.AllowAllMethods,
				AllowAllOrigins:  instance.Spec.CORSHeaders.AllowAllOrigins,
				AllowAllHeaders:  instance.Spec.CORSHeaders.AllowAllHeaders,
				AllowCredentials: instance.Spec.CORSHeaders.AllowCredentials,
				MaxAge:           instance.Spec.CORSHeaders.MaxAge,
			}

			accessApplication.CorsHeaders = &corsHeaders
		}

		if err := createAccessApplication(cf, &accessApplication); err != nil {
			return nil, errors.Wrap(err, "failed to create access application")
		}

		return &accessApplication, nil
	}

	existingApplication, err := getExistingAccessApplicationFromID(instance, zone, cf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get existing application")
	}

	if existingApplication == nil {
		// this happens when we have a status resource with an id but it was deleted from cloudflare
		// TODO we need to recreate the app and update the status resource here
		return nil, errors.New("application not found, maybe it was deleted.  update the status.applicationID to remove it from this resource and reconcile again")
	}

	hasChanges := false

	toUpdate := cloudflare.AccessApplication{
		ID:     existingApplication.ID,
		Name:   instance.Spec.Name,
		Domain: instance.Spec.Domain,
	}

	if instance.Spec.Name != existingApplication.Name {
		hasChanges = true
	}
	if instance.Spec.Domain != existingApplication.Domain {
		hasChanges = true
	}

	if instance.Spec.SessionDuration != "" {
		if instance.Spec.SessionDuration != existingApplication.SessionDuration {
			hasChanges = true
			toUpdate.SessionDuration = instance.Spec.SessionDuration
		}
	}
	if instance.Spec.AllowedIdPs != nil {
		sort.Strings(instance.Spec.AllowedIdPs)
		sort.Strings(existingApplication.AllowedIdps)

		idpsChanged := len(instance.Spec.AllowedIdPs) != len(existingApplication.AllowedIdps)

		if !idpsChanged {
			for i, v := range instance.Spec.AllowedIdPs {
				if v != existingApplication.AllowedIdps[i] {
					idpsChanged = true
				}
			}
		}

		if idpsChanged {
			hasChanges = true
			toUpdate.AllowedIdps = instance.Spec.AllowedIdPs
		}
	}
	if instance.Spec.AutoRedirectToIdentity != nil {
		if existingApplication.AutoRedirectToIdentity != *instance.Spec.AutoRedirectToIdentity {
			hasChanges = true
			toUpdate.AutoRedirectToIdentity = *instance.Spec.AutoRedirectToIdentity
		}
	}
	if instance.Spec.CORSHeaders != nil {
		updatedCORSHeaders, err := diffCORSHeaders(existingApplication.CorsHeaders, instance.Spec.CORSHeaders)
		if err != nil {
			return nil, errors.Wrap(err, "failed to diff cors headers")
		}
		if updatedCORSHeaders != nil {
			hasChanges = true
			toUpdate.CorsHeaders = updatedCORSHeaders
		}
	}

	if hasChanges {
		if err := updateAccessApplication(cf, &toUpdate); err != nil {
			return nil, errors.Wrap(err, "failed to update access application")
		}

		existingApplication = &toUpdate
	}

	if len(instance.Spec.AccessPolicies) > 0 {
		// TODO
		existingAccessPolicies := []cloudflare.AccessPolicy{} // TODO

		toCreate, toUpdate, toDelete, err := diffAccessPolicies(existingAccessPolicies, instance.Spec.AccessPolicies)
		if err != nil {
			return nil, errors.Wrap(err, "failed to diff access policies")
		}

		for _, ap := range toCreate {
			_, err := cf.CreateAccessPolicy(cf.AccountID, existingApplication.ID, ap)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create access policy")
			}
		}
		for _, ap := range toUpdate {
			_, err := cf.UpdateAccessPolicy(cf.AccountID, existingApplication.ID, ap)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update access policy")
			}
		}
		for _, ap := range toDelete {
			err := cf.DeleteAccessPolicy(cf.AccountID, existingApplication.ID, ap.ID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to delete access policy")
			}
		}
	}

	return existingApplication, nil
}

// getExistingAccessApplicationFromID will look at the status subresource and if there's an ID
// will get that. if not, it will a
func getExistingAccessApplicationFromID(instance crdsv1alpha1.AccessApplication, zone *crdsv1alpha1.Zone, cf *cloudflare.API) (*cloudflare.AccessApplication, error) {
	accessApplication, err := cf.AccessApplication(cf.AccountID, instance.Status.ApplicationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get accessapplication from cf")
	}

	return &accessApplication, nil
}

func findExistingAccessApplication(instance crdsv1alpha1.AccessApplication, zone *crdsv1alpha1.Zone, cf *cloudflare.API) (*cloudflare.AccessApplication, error) {
	currentPage := cloudflare.PaginationOptions{
		PerPage: 20,
		Page:    0,
	}

	for {
		accessApplications, resultInfo, err := cf.AccessApplications(cf.AccountID, currentPage)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get page of accessapplications")
		}

		for _, accessApplication := range accessApplications {
			if accessApplication.Domain == instance.Spec.Domain && accessApplication.Name == instance.Spec.Name {
				return &accessApplication, nil
			}
		}

		currentPage.Page++

		if currentPage.Page >= resultInfo.TotalPages {
			return nil, nil
		}
	}
}

func createAccessApplication(cf *cloudflare.API, application *cloudflare.AccessApplication) error {
	createdApplication, err := cf.CreateAccessApplication(cf.AccountID, *application)
	if err != nil {
		return errors.Wrap(err, "failed to create access application")
	}

	application = &createdApplication
	return nil
}

func updateAccessApplication(cf *cloudflare.API, application *cloudflare.AccessApplication) error {
	updatedApplication, err := cf.UpdateAccessApplication(cf.AccountID, *application)
	if err != nil {
		return errors.Wrap(err, "failed to update access application")
	}

	application = &updatedApplication
	return nil
}
