package accessapplication

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/controller/shared"
)

func diffAccessPolicies(existings []cloudflare.AccessPolicy, desireds []crdsv1alpha1.AccessPolicy) ([]cloudflare.AccessPolicy, []cloudflare.AccessPolicy, []cloudflare.AccessPolicy, error) {
	toCreate := []cloudflare.AccessPolicy{}
	toUpdate := []cloudflare.AccessPolicy{}
	toDelete := []cloudflare.AccessPolicy{}

	if len(existings) == 0 && len(desireds) == 0 {
		return toCreate, toUpdate, toDelete, nil
	}

	for _, existing := range existings {
		found := false
		for _, desired := range desireds {
			// TODO we should store these by ID in the status field
			// with this implementation, renaming would delete and recreate
			if desired.Name == existing.Name {
				found = true

				updated, err := diffAccessPolicy(existing, desired)
				if err != nil {
					return nil, nil, nil, errors.Wrap(err, "failed to diff policies")
				}

				if updated != nil {
					toUpdate = append(toUpdate, *updated)
				}

				goto Found
			}
		}

	Found:
		if !found {
			toDelete = append(toDelete, existing)
		}
	}

	for _, desired := range desireds {
		found := false
		for _, existing := range existings {
			if existing.Name == desired.Name {
				found = true
				goto Found2
			}
		}

	Found2:
		if !found {
			create := cloudflare.AccessPolicy{
				Name:     desired.Name,
				Decision: desired.Decision,
				Include:  shared.StringArrayToInterfaceArray(desired.Include),
				Exclude:  shared.StringArrayToInterfaceArray(desired.Exclude),
				Require:  shared.StringArrayToInterfaceArray(desired.Require),
			}

			if desired.Precedence != nil {
				create.Precedence = *desired.Precedence
			}

			toCreate = append(toCreate, create)
		}
	}
	return toCreate, toUpdate, toDelete, nil
}

// diffAccessPolicy will diff existing to desired.
// if there are diffs, it will return not-nil in the first response param
func diffAccessPolicy(existing cloudflare.AccessPolicy, desired crdsv1alpha1.AccessPolicy) (*cloudflare.AccessPolicy, error) {
	hasChanged := false

	if existing.Name != desired.Name {
		hasChanged = true
		existing.Name = desired.Name
	}

	if existing.Decision != desired.Decision {
		hasChanged = true
		existing.Decision = desired.Decision
	}

	if desired.Precedence != nil {
		if existing.Precedence != *desired.Precedence {
			hasChanged = true
			existing.Precedence = *desired.Precedence
		}
	}

	if !shared.StringSlicesMatch(shared.InterfaceArrayToStringArray(existing.Include), desired.Include) {
		hasChanged = true
		existing.Include = shared.StringArrayToInterfaceArray(desired.Include)
	}

	if !shared.StringSlicesMatch(shared.InterfaceArrayToStringArray(existing.Exclude), desired.Exclude) {
		hasChanged = true
		existing.Exclude = shared.StringArrayToInterfaceArray(desired.Exclude)
	}

	if !shared.StringSlicesMatch(shared.InterfaceArrayToStringArray(existing.Require), desired.Require) {
		hasChanged = true
		existing.Require = shared.StringArrayToInterfaceArray(desired.Require)
	}

	if !hasChanged {
		return nil, nil
	}

	return &existing, nil
}
