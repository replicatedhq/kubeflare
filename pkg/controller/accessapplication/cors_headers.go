package accessapplication

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/controller/shared"
)

// diffCORSHeaders will diff existing to desired.
// if there are diffs, it will return not-nil in the first response param
func diffCORSHeaders(existing *cloudflare.AccessApplicationCorsHeaders, desired *crdsv1alpha1.CORSHeader) (*cloudflare.AccessApplicationCorsHeaders, error) {
	if existing == nil && desired == nil {
		return nil, nil
	}

	if existing == nil {
		// copy all of the desired
		updated := &cloudflare.AccessApplicationCorsHeaders{
			AllowedMethods:   desired.AllowedMethods,
			AllowedOrigins:   desired.AllowedOrigins,
			AllowedHeaders:   desired.AllowedHeaders,
			AllowAllMethods:  desired.AllowAllMethods,
			AllowAllOrigins:  desired.AllowAllOrigins,
			AllowAllHeaders:  desired.AllowAllHeaders,
			AllowCredentials: desired.AllowCredentials,
			MaxAge:           desired.MaxAge,
		}

		return updated, nil
	}

	if desired == nil {
		// TODO, delete?
		return nil, errors.New("deleting cors headers is not yet implemented")
	}

	wasUpdated := false
	updated := cloudflare.AccessApplicationCorsHeaders{
		AllowedMethods:   desired.AllowedMethods,
		AllowedOrigins:   desired.AllowedOrigins,
		AllowedHeaders:   desired.AllowedHeaders,
		AllowAllMethods:  desired.AllowAllMethods,
		AllowAllOrigins:  desired.AllowAllOrigins,
		AllowAllHeaders:  desired.AllowAllHeaders,
		AllowCredentials: desired.AllowCredentials,
		MaxAge:           desired.MaxAge,
	}

	if !shared.StringSlicesMatch(existing.AllowedMethods, desired.AllowedMethods) {
		wasUpdated = true
	}
	if !shared.StringSlicesMatch(existing.AllowedOrigins, desired.AllowedOrigins) {
		wasUpdated = true
	}

	if !shared.StringSlicesMatch(existing.AllowedHeaders, desired.AllowedHeaders) {
		wasUpdated = true
	}

	if existing.AllowAllMethods != desired.AllowAllMethods {
		wasUpdated = true
	}

	if existing.AllowAllOrigins != desired.AllowAllOrigins {
		wasUpdated = true
	}

	if existing.AllowAllHeaders != desired.AllowAllHeaders {
		wasUpdated = true
	}

	if existing.MaxAge != desired.MaxAge {
		wasUpdated = true
	}

	if !wasUpdated {
		return nil, nil
	}

	return &updated, nil
}
