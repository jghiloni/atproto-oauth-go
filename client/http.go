package client

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jghiloni/atproto-oauth-go/retry"
)

var (
	ErrInvalidScheme  = errors.New("invalid scheme. only https and http with localhost supported")
	ErrPrivateAddress = errors.New("private addresses not allowed")
)

const DefaultClientTimeout = 10 * time.Second

type hardenedRoundTripper struct {
	delegate http.RoundTripper
}

func (h *hardenedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := h.validateRequest(req); err != nil {
		return nil, err
	}

	strategy := &retry.LinearRetryStrategy{
		MaxRetries: 10,
		Backoff:    2 * time.Second,
		MaxBackoff: 10 * time.Second,
	}

	return retry.DoWithRetry(func() (*http.Response, error) {
		r, e := h.delegate.RoundTrip(req)
		if e == nil && r.StatusCode >= 500 {
			return nil, retry.ErrRetry
		}
		return r, e
	}, strategy)
}

func (h *hardenedRoundTripper) validateRequest(req *http.Request) error {
	scheme := strings.ToLower(req.URL.Scheme)
	if scheme != "https" {
		if scheme == "http" {
			if !strings.EqualFold("localhost", req.URL.Host) {
				return fmt.Errorf("%w. got %s", ErrInvalidScheme, scheme)
			}
		}
	}

	ip := net.ParseIP(req.URL.Hostname())
	ips := []net.IPAddr{}

	var err error
	if ip == nil {
		ips, err = net.DefaultResolver.LookupIPAddr(req.Context(), req.URL.Hostname())
		if err != nil {
			return err
		}
	} else {
		ips = append(ips, net.IPAddr{
			IP: ip,
		})
	}

	for _, ipAddr := range ips {
		if ipAddr.IP.IsPrivate() {
			return ErrPrivateAddress
		}
	}

	return nil
}

func HardenHTTPClient(hc *http.Client) *http.Client {
	if hc == nil {
		hc = http.DefaultClient
	}

	if hc.Transport == nil {
		hc.Transport = http.DefaultTransport
	}

	timeout := hc.Timeout
	if timeout == 0 {
		timeout = DefaultClientTimeout
	}

	return &http.Client{
		Transport: &hardenedRoundTripper{
			delegate: hc.Transport,
		},
		CheckRedirect: hc.CheckRedirect,
		Jar:           hc.Jar,
		Timeout:       timeout,
	}
}
