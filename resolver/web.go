package resolver

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type webResolver struct {
	hc *http.Client
}

func NewWebResolver(client ...*http.Client) DocumentResolver {
	hc := defaultHC
	if len(client) > 0 {
		hc = client[0]
	}

	return &webResolver{
		hc,
	}
}

func (*webResolver) SupportedDIDTypes() []string {
	return []string{"web"}
}

func (w *webResolver) ResolveDIDDocumentURL(ctx context.Context, didResolver DIDResolver, handleOrDID string) (string, error) {
	did, err := resolveIdentifierToDID(ctx, didResolver, handleOrDID)
	if err != nil {
		return "", fmt.Errorf("could not determine did for identifier %s: %w", handleOrDID, err)
	}

	if err := verifyDIDTypeForResolver(did, w); err != nil {
		return "", err
	}

	baseURL := strings.ReplaceAll(did[8:], ":", "/")
	baseURL, err = url.PathUnescape(baseURL)
	if err != nil {
		return "", ErrMalformedDID
	}

	return fmt.Sprintf("https://%s/.well-known/did.json", baseURL), nil
}
