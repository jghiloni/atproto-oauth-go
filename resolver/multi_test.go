package resolver_test

import (
	"context"
	"testing"

	"github.com/jghiloni/atproto-oauth-go/resolver"
	. "github.com/onsi/gomega"
)

func TestMultiResolver(t *testing.T) {
	tests := []struct {
		handle        string
		expectedDID   string
		expectedError bool
	}{
		{"@jaygles.bsky.social", "did:plc:e2fun4xcfwtcrqfdwhfnghxk", false},
		{"watchedsky.social", "did:plc:hvjfuy2w6zqu6abmpkwcpulc", false},
		{"&invalid.invalid", "", true},
	}

	RegisterTestingT(t)
	for _, tc := range tests {
		actualDID, err := resolver.DefaultDIDResolver.ResolveHandleToDID(context.Background(), tc.handle)
		if tc.expectedError {
			Expect(err).Should(HaveOccurred())
		} else {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actualDID).To(Equal(tc.expectedDID))
		}
	}
}
