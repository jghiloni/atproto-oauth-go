package resolver_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jghiloni/atproto-oauth-go/resolver"
	. "github.com/onsi/gomega"
)

func TestHTTPSResolver(t *testing.T) {
	tests := []struct {
		clients          []*http.Client
		handle           string
		expectedresolver string
		errorExpected    bool
	}{
		{[]*http.Client{}, "@jaygles.bsky.social", "did:plc:e2fun4xcfwtcrqfdwhfnghxk", false},
		{[]*http.Client{http.DefaultClient}, "jaygles.bsky.social", "did:plc:e2fun4xcfwtcrqfdwhfnghxk", false},
		{nil, "watchedsky.social", "", true},
		{nil, "invalid.invalid", "", true},
	}

	RegisterTestingT(t)
	for _, testCase := range tests {
		r := resolver.NewHTTPSResolver(testCase.clients...)
		resolver, err := r.ResolveHandleToDID(context.Background(), testCase.handle)
		if testCase.errorExpected {
			Expect(err).Should(HaveOccurred())
		} else {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resolver).To(Equal(testCase.expectedresolver))
		}
	}
}
