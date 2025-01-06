package resolver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jghiloni/atproto-oauth-go/client"
)

type httpsResolver struct {
	hc *http.Client
}

var defaultHC = client.HardenHTTPClient(http.DefaultClient)

// NewHTTPSResolver returns an AtprotoHandleResolver that resolves DIDs from Handles by
// accessing https://<handle>/.well-known/atproto-did. If you pass multiple *http.Clients,
// it only uses the first one; this is just a workaround to have a no-args and 1 arg
// method. If no client is provided, a default client will be used that behaves like
// http.DefaultClient with up to 10 backoffs, using a linear backoff strategy
func NewHTTPSResolver(clients ...*http.Client) DIDResolver {
	hc := defaultHC
	if len(clients) > 0 {
		hc = clients[0]
	}

	return &httpsResolver{
		hc: hc,
	}
}

// ResolveHandleToDID will take a handle like @jaygles.bsky.social (leading @ optional) and look for its
// did at https://jaygles.bsky.social/.well-known/atproto-did, as described in the AT Protocol documentation.
func (h *httpsResolver) ResolveHandleToDID(ctx context.Context, handle string) (did string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/.well-known/atproto-did", strings.TrimPrefix(handle, "@")), nil)
	if err != nil {
		return "", fmt.Errorf("could not generate request: %w", err)
	}

	resp, err := h.hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("error requesting did: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected status 200 OK, got %s", resp.Status)
	}

	ct := resp.Header.Get("content-type")
	if !strings.HasPrefix(ct, "text/plain") {
		return "", fmt.Errorf("expected Content-Type text/plain, got %s", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	return string(body), nil
}
