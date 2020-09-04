package accessapplication

import (
	"testing"

	"github.com/cloudflare/cloudflare-go"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func Test_diffCORSHeaders(t *testing.T) {
	tests := []struct {
		name     string
		existing *cloudflare.AccessApplicationCorsHeaders
		desired  *crdsv1alpha1.CORSHeader
		expected *cloudflare.AccessApplicationCorsHeaders
	}{
		{
			name:     "both are nil",
			existing: nil,
			desired:  nil,
			expected: nil,
		},
		{
			name:     "no existing, has desired",
			existing: nil,
			desired: &crdsv1alpha1.CORSHeader{
				AllowedMethods:   []string{"a"},
				AllowedOrigins:   []string{"b", "c"},
				AllowedHeaders:   []string{},
				AllowAllMethods:  false,
				AllowAllOrigins:  true,
				AllowAllHeaders:  false,
				AllowCredentials: true,
				MaxAge:           9999,
			},
			expected: &cloudflare.AccessApplicationCorsHeaders{
				AllowedMethods:   []string{"a"},
				AllowedOrigins:   []string{"b", "c"},
				AllowedHeaders:   []string{},
				AllowAllMethods:  false,
				AllowAllOrigins:  true,
				AllowAllHeaders:  false,
				AllowCredentials: true,
				MaxAge:           9999,
			},
		},
		{
			name: "change 1 field",
			existing: &cloudflare.AccessApplicationCorsHeaders{
				AllowedMethods:   []string{"a"},
				AllowedOrigins:   []string{"b", "c"},
				AllowedHeaders:   []string{},
				AllowAllMethods:  false,
				AllowAllOrigins:  true,
				AllowAllHeaders:  false,
				AllowCredentials: true,
				MaxAge:           9999,
			},
			desired: &crdsv1alpha1.CORSHeader{
				AllowedMethods:   []string{"a"},
				AllowedOrigins:   []string{"b", "c'"},
				AllowedHeaders:   []string{},
				AllowAllMethods:  false,
				AllowAllOrigins:  true,
				AllowAllHeaders:  false,
				AllowCredentials: true,
				MaxAge:           9999,
			},
			expected: &cloudflare.AccessApplicationCorsHeaders{
				AllowedMethods:   []string{"a"},
				AllowedOrigins:   []string{"b", "c'"},
				AllowedHeaders:   []string{},
				AllowAllMethods:  false,
				AllowAllOrigins:  true,
				AllowAllHeaders:  false,
				AllowCredentials: true,
				MaxAge:           9999,
			},
		},
	}

	for _, test := range tests {
		req := require.New(t)
		actual, err := diffCORSHeaders(test.existing, test.desired)
		req.NoError(err)
		assert.DeepEqual(t, test.expected, actual)
	}
}
