package resolver

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"
)

// taken from github.com/bluesky-social/indigo/atproto/identity to avoid the dependency

type DocVerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

type DocService struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

type Document struct {
	DID                string                  `json:"id"`
	AlsoKnownAs        []string                `json:"alsoKnownAs,omitempty"`
	VerificationMethod []DocVerificationMethod `json:"verificationMethod,omitempty"`
	Service            []DocService            `json:"service,omitempty"`
}

const (
	pdsID   = "#atproto_pds"
	pdsType = "AtprotoPersonalDataServer"
)

var (
	ErrInvalidDIDType = errors.New("invalid DID type for this resolver")
	ErrMalformedDID   = errors.New("malformed DID")
	ErrNoPDSFound     = fmt.Errorf("no service of type %s found in document", pdsType)
)

func resolveIdentifierToDID(ctx context.Context, resolver DIDResolver, identifier string) (string, error) {
	if strings.HasPrefix(identifier, "did:") {
		return identifier, nil
	}

	if resolver == nil {
		resolver = DefaultDIDResolver
	}

	return resolver.ResolveHandleToDID(ctx, identifier)
}

func verifyDIDTypeForResolver(did string, resolver DocumentResolver) error {
	if resolver == nil {
		return errors.New("resolver cannot be nil")
	}

	parts := strings.SplitN(did, ":", 3)
	if len(parts) != 3 {
		return fmt.Errorf("%w: not enough parts", ErrMalformedDID)
	}

	if parts[0] != "did" {
		return fmt.Errorf("%w: first segment must be 'did'", ErrMalformedDID)
	}

	if slices.Contains(resolver.SupportedDIDTypes(), parts[1]) {
		return fmt.Errorf("%w: this resolver supports DID types %q but got %q", ErrInvalidDIDType, strings.Join(resolver.SupportedDIDTypes(), ","), parts[1])
	}

	return nil
}

func GetPDSFromDIDDocument(doc Document, strict bool) (string, error) {
	for _, svc := range doc.Service {
		if svc.ID == pdsID {
			if strict {
				if svc.Type != pdsType {
					return "", fmt.Errorf("expected service type %q, but got %q", pdsType, svc.Type)
				}

				_, err := url.Parse(svc.ServiceEndpoint)
				if err != nil {
					return "", fmt.Errorf("service endpoint is an invalid URL: %w", err)
				}
			}

			return svc.ServiceEndpoint, nil
		}
	}

	return "", ErrNoPDSFound
}
