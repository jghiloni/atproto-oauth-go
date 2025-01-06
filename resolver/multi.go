package resolver

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
)

type multiDIDResolver struct {
	rs []DIDResolver
	t  time.Duration
}

var ErrTimedOut = errors.New("resolution timed out")

func NewParallelResolver(timeout time.Duration, resolvers ...DIDResolver) DIDResolver {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &multiDIDResolver{
		rs: resolvers,
		t:  timeout,
	}
}

func (m *multiDIDResolver) ResolveHandleToDID(ctx context.Context, handle string) (string, error) {
	subCtx, cancel := context.WithCancel(ctx)
	errs := make(chan error, len(m.rs))
	answer := make(chan string, 1)
	defer func() {
		_ = recover()
		close(errs)
		cancel()
	}()

	for _, r := range m.rs {
		go m.resolveSingle(subCtx, handle, r, answer, errs)
	}

	go func() {
		defer func() {
			_ = recover()
		}()

		errSlice := []error{}
		for e := range errs {
			errSlice = append(errSlice, e)
			if len(m.rs) == len(errSlice) {
				close(answer)
				errs <- errors.Join(errSlice...)
				return
			}
		}

	}()

	select {
	case txt, ok := <-answer:
		if ok {
			close(answer)
			return txt, nil
		}

		return "", <-errs
	case <-time.After(m.t):
		return "", ErrTimedOut
	}
}

func (m *multiDIDResolver) resolveSingle(ctx context.Context, handle string, r DIDResolver, answer chan string, errs chan error) {
	defer func() {
		_ = recover()
	}()

	txt, err := r.ResolveHandleToDID(ctx, handle)
	if err != nil {
		errs <- err
		return
	}

	answer <- txt
}

type routingDocumentResolver map[string]DocumentResolver

func (r routingDocumentResolver) SupportedDIDTypes() []string {
	return slices.AppendSeq([]string{}, maps.Keys(r))
}

func (r routingDocumentResolver) ResolveDIDDocumentURL(ctx context.Context, d DIDResolver, identity string) (string, error) {
	did, err := resolveIdentifierToDID(ctx, d, identity)
	if err != nil {
		return did, err
	}

	if err := verifyDIDTypeForResolver(did, r); err != nil {
		return "", err
	}

	didType := strings.SplitN(did, ":", 3)[1]
	docResolver, ok := r[didType]
	if !ok {
		return "", fmt.Errorf("no resolver found for did %s", did)
	}

	return docResolver.ResolveDIDDocumentURL(ctx, d, did)
}
