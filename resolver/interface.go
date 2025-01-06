package resolver

import "context"

type DIDResolver interface {
	ResolveHandleToDID(ctx context.Context, handle string) (did string, err error)
}

type DocumentResolver interface {
	SupportedDIDTypes() []string
	ResolveDIDDocumentURL(ctx context.Context, didResolver DIDResolver, handleOrDID string) (string, error)
}
