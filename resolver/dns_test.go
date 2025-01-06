package resolver_test

import (
	"context"
	"testing"

	"github.com/jghiloni/atproto-oauth-go/resolver"
	. "github.com/onsi/gomega"
)

func TestDNSResolver(t *testing.T) {
	tests := []struct {
		nameservers   []string
		handle        string
		expectedDid   string
		errorExpected bool
	}{
		{[]string{}, "@jaygles.bsky.social", "", true},
		{[]string{}, "@watchedsky.social", "did:plc:hvjfuy2w6zqu6abmpkwcpulc", false},
		{[]string{"8.8.8.8"}, "watchedsky.social", "did:plc:hvjfuy2w6zqu6abmpkwcpulc", false},
		{[]string{"1.1.1.1", "8.8.8.8"}, "invalid.invalid", "", true},
	}

	RegisterTestingT(t)
	for _, testCase := range tests {
		r := resolver.NewDNSResolver(testCase.nameservers...)
		did, err := r.ResolveHandleToDID(context.Background(), testCase.handle)
		if testCase.errorExpected {
			Expect(err).Should(HaveOccurred())
		} else {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(did).To(Equal(testCase.expectedDid))
		}
	}
}
