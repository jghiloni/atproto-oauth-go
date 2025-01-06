package resolver

import "time"

var DefaultDIDResolver = NewParallelResolver(20*time.Second, NewDNSResolver(), NewHTTPSResolver())

var DefaultDocumentResolver = routingDocumentResolver{
	"plc": NewPLCResolver(),
	"web": NewWebResolver(),
}

var DefaultDocumentFetcher = NewDocumentFetcher(DefaultDIDResolver, DefaultDocumentResolver, defaultHC)
