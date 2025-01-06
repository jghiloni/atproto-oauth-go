package resolver

import (
	"context"
	"fmt"
	"net/http"
)

type plcResolver struct {
	hc *http.Client
}

func NewPLCResolver(client ...*http.Client) DocumentResolver {
	hc := defaultHC
	if len(client) > 0 {
		hc = client[0]
	}
	return &plcResolver{
		hc,
	}
}

const plcBaseURL = "https://plc.directory"

func (*plcResolver) SupportedDIDTypes() []string {
	return []string{"plc"}
}

func (p *plcResolver) ResolveDIDDocumentURL(ctx context.Context, didResolver DIDResolver, handleOrDid string) (didURL string, err error) {
	did, err := resolveIdentifierToDID(ctx, didResolver, handleOrDid)
	if err != nil {
		return "", fmt.Errorf("could not determine did for identifier %s: %w", handleOrDid, err)
	}

	if err := verifyDIDTypeForResolver(did, p); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", plcBaseURL, did), nil
}
