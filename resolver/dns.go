package resolver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type dnsResolver struct {
	nameServers []string
}

// TXTRecordPrefix is used to determine if a TXT record is appropriate
const TXTRecordPrefix = "did="

// ErrDNSResolutionFailed is returned if there is no appropriate _atproto TXT record
var ErrDNSResolutionFailed = &net.DNSError{
	IsTemporary: false,
	IsNotFound:  true,
	Name:        "E_NOTFOUND",
}

var errTemporary = errors.New("temporary")

// NewDNSResolver returns an AtprotoHandleResolver that finds the DID for
// a handle by looking for a DNS TXT record at _atproto.<handle>. Pass in a
// list of nameservers to use, and they will be tried in order until the
// query succeeds or there is a non-temporary error. If no nameservers are
// specified, the default nameservers for the host will be used.
func NewDNSResolver(nameservers ...string) DIDResolver {
	return &dnsResolver{
		nameServers: nameservers,
	}
}

func dialer(address string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, _ string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: time.Second,
		}

		return d.DialContext(ctx, network, address)
	}
}

// ResolveHandleToDID attempts to find the DID associated with a given handle by searching
// for an associated DNS TXT record at _atproto.<handle>.
func (d *dnsResolver) ResolveHandleToDID(ctx context.Context, handle string) (string, error) {
	resolver := &net.Resolver{
		PreferGo: true,
	}

	if len(d.nameServers) == 0 {
		return d.resolveWithResolver(ctx, net.DefaultResolver, handle)
	}

	for _, addr := range d.nameServers {
		resolver.Dial = dialer(fmt.Sprintf("%s:53", addr))
		txt, err := d.resolveWithResolver(ctx, resolver, handle)
		if err != nil {
			if errors.Is(err, errTemporary) {
				continue
			}

			return "", err
		}

		return txt, nil
	}

	return "", ErrDNSResolutionFailed
}

func (d *dnsResolver) resolveWithResolver(ctx context.Context, resolver *net.Resolver, handle string) (string, error) {
	txt, err := resolver.LookupTXT(ctx, fmt.Sprintf("_atproto.%s", strings.TrimPrefix(handle, "@")))
	if err != nil {
		de := new(net.DNSError)
		if errors.As(err, &de) {
			if de.Temporary() {
				return "", errTemporary
			}
		}

		return "", err
	}

	for _, txtRecord := range txt {
		if strings.HasPrefix(txtRecord, TXTRecordPrefix) {
			return strings.TrimPrefix(txtRecord, TXTRecordPrefix), nil
		}
	}

	return "", ErrDNSResolutionFailed
}
