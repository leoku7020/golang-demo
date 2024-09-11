package servkit

import (
	"context"
	"net/textproto"
	"strconv"
	"strings"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

const (
	SpecifiedHeaderPrefix = "spec-"

	HTTPMethod        = SpecifiedHeaderPrefix + "http-method"
	HTTPRealIP        = SpecifiedHeaderPrefix + "http-real-ip"
	HTTPRequestURI    = SpecifiedHeaderPrefix + "http-request-uri"
	HTTPRoute         = SpecifiedHeaderPrefix + "http-route"
	HTTPUserAgent     = SpecifiedHeaderPrefix + "http-user-agent"
	HTTPCode          = SpecifiedHeaderPrefix + "http-code"
	HTTPAuthorization = SpecifiedHeaderPrefix + "authorization"
)

var (
	retHTTPCode = renderHeaderKey(gwruntime.MetadataHeaderPrefix + HTTPCode)
)

type MD struct {
	md metadata.MD
}

func GetMetadata(ctx context.Context) (*MD, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, false
	}

	return &MD{md}, true
}

func (md *MD) Method() string {
	return md.getSpecKey(HTTPMethod)
}

func (md *MD) RealIP() string {
	return md.getSpecKey(HTTPRealIP)
}

func (md *MD) RequestURI() string {
	return md.getSpecKey(HTTPRequestURI)
}

func (md *MD) Route() string {
	return md.getSpecKey(HTTPRoute)
}

func (md *MD) UserAgent() string {
	return md.getSpecKey(HTTPUserAgent)
}

func (md *MD) Authorization() string {
	return md.getSpecKey(HTTPAuthorization)
}

func (md *MD) Header(key string) string {
	return md.getPartalKey(key)
}

func (md *MD) getSpecKey(key string) string {
	k := strings.ToLower(key)
	vals, ok := md.md[k]
	if !ok || len(vals) == 0 {
		return ""
	}

	return vals[len(vals)-1]
}

func (md *MD) getPartalKey(key string) string {
	if val := md.getSpecKey(SpecifiedHeaderPrefix + key); val != "" {
		return val
	}

	if val := md.getSpecKey(gwruntime.MetadataPrefix + key); val != "" {
		return val
	}

	return ""
}

func renderHeaderKey(key string) string {
	return textproto.CanonicalMIMEHeaderKey(key)
}

func extractStatusCode(ctx context.Context) (int, error) {
	md, ok := gwruntime.ServerMetadataFromContext(ctx)
	if !ok {
		return 0, nil
	}

	if vals := md.HeaderMD.Get(HTTPCode); len(vals) > 0 {
		code, err := strconv.Atoi(vals[len(vals)-1])
		if err != nil {
			return 0, err
		}

		return code, nil
	}

	return 0, nil
}
