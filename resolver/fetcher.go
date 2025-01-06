package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DocumentFetcher struct {
	didResolver DIDResolver
	docResolver DocumentResolver
	hc          *http.Client
}

func NewDocumentFetcher(didResolver DIDResolver, docResolver DocumentResolver, hc *http.Client) *DocumentFetcher {
	return &DocumentFetcher{
		didResolver, docResolver, hc,
	}
}

func (d *DocumentFetcher) GetDocumentForIdentifier(ctx context.Context, handleOrDID string) (*Document, error) {
	url, err := d.docResolver.ResolveDIDDocumentURL(ctx, d.didResolver, handleOrDID)
	if err != nil {
		return nil, fmt.Errorf("could not get document URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not build request: %w", err)
	}

	resp, err := d.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get a response from url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status 200 OK, got %s", resp.Status)
	}

	var doc Document
	if err = json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("could not parse DID document: %w", err)
	}

	return &doc, nil
}
