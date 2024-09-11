package servkit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tomasen/realip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"demo/internal/router/middleware"
	"demo/pkg/bootkit"
	"demo/pkg/envkit"
	"demo/pkg/httpkit"
	"demo/pkg/logger"
)

var (
	// ErrInvalidStatusCode indicated not existing code
	ErrInvalidStatusCode = errors.New("invalid status code")
)

type RegisterHandlerFunc func(ctx context.Context, mux *gwruntime.ServeMux, conn *grpc.ClientConn) error

func RunGrpcGateway(
	ctx context.Context,
	gwAddr, grpcAddr string,
	shutdown bootkit.ShutdownFunc,
	registerHandlers RegisterHandlerFunc,
	options ...ServOptions,
) error {
	o := applyServOptions(options...)

	o.grpcDialOpts = append(o.grpcDialOpts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.DialContext(ctx, grpcAddr, o.grpcDialOpts...)
	if err != nil {
		logger.Ctx(ctx).Error("grpc.DialContext failed",
			logger.WithError(err),
			logger.WithField("grpcAddr", grpcAddr),
		)
		return err
	}

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/health", healthGRPCServer(conn))
	// k8s API health endpoints
	// ref: https://kubernetes.io/docs/reference/using-api/health-checks/
	httpMux.HandleFunc("/healthz", healthzGRPCServer(conn))

	o.gwServMuxOpts = append(o.gwServMuxOpts, enrichMetaData())
	gwMux := gwruntime.NewServeMux(o.gwServMuxOpts...)
	if err := registerHandlers(ctx, gwMux, conn); err != nil {
		logger.Ctx(ctx).Error("registerHandler failed", logger.WithError(err))
		return err
	}

	// don't allow CORS in Production
	if envkit.Namespace() != envkit.EnvProduction {
		o.middlewares = append(o.middlewares, middleware.AllowCORS)
	}

	// check if content-encoding is gzip, and decompress it
	o.middlewares = append(o.middlewares, middleware.GZipDecompressor, middleware.PayloadMiddleware)

	chain := httpkit.NewChain(o.middlewares...)

	httpMux.Handle("/", gwMux)

	s := &http.Server{
		Addr:    gwAddr,
		Handler: chain.Then(httpMux),
	}

	shutdown(func() error {
		// close gRPC gateway
		logger.Info("Shutting down gRPC gateway ...")
		if err := s.Shutdown(context.Background()); err != nil {
			return err
		}

		// close connection
		return conn.Close()
	})

	logger.Ctx(ctx).Info(fmt.Sprintf("Starting gRPC gateway listening at %s", gwAddr))

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logger.Ctx(ctx).Error("Failed to listen and serve", logger.WithError(err))
		return err
	}

	return nil
}

// healthGRPCServer returns a simple health handler which returns {"ok":true}.
func healthGRPCServer(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf(`{"ok":false,"reason":"%s"}`, s), http.StatusBadGateway)
			return
		}

		fmt.Fprintln(w, `{"ok":true}`)
	}
}

// healthzGRPCServer returns a simple health handler which returns ok.
func healthzGRPCServer(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}

// openAPIServer returns OpenAPI specification files located under "/openapiv2/"
func openAPIServer(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
			logger.Error(fmt.Sprintf("Not Found: %s", r.URL.Path))
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(r.URL.Path, "/openapiv2/")
		p = path.Join(dir, p)
		logger.Info(fmt.Sprintf("Serving %s", r.URL.Path), logger.WithField("path", p))
		http.ServeFile(w, r, p)
	}
}

func SetFormULREncodeMarhslalerOptions() gwruntime.ServeMuxOption {
	return gwruntime.WithMarshalerOption("application/x-www-form-urlencoded", UrlEncodeMarshal{})
}

func SetFileEncodeMarhslalerOptions() gwruntime.ServeMuxOption {
	return gwruntime.WithMarshalerOption("multipart/form-data", FileMarshal{})
}

func SetPdfEncodeMarhslalerOptions() gwruntime.ServeMuxOption {
	//return gwruntime.WithMarshalerOption(gwruntime.MIMEWildcard, &gwruntime.HTTPBodyMarshaler{})
	return gwruntime.WithMarshalerOption("application/json", CustomJsonMarshal{})
}

func SetTextEncodeMarhslalerOptions() gwruntime.ServeMuxOption {
	return gwruntime.WithMarshalerOption("text/plain", TextPlainMarshal{})
}

// SetDefaultMarshalerOptions sets Marshaler to add default values into unpopulated fields during marshaling,
// and ignore unknow fields during unmarshaling.
func SetDefaultMarshalerOptions() gwruntime.ServeMuxOption {
	// JSONPb is a Marshaler which marshals/unmarshals into/from JSON
	// with the "google.golang.org/protobuf/encoding/protojson" marshaler.
	// It supports the full functionality of protobuf unlike JSONBuiltin.
	return gwruntime.WithMarshalerOption(gwruntime.MIMEWildcard, &gwruntime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			// TL;DR keep unpopulated fields with the default values
			// EmitUnpopulated specifies whether to emit unpopulated fields. It does not
			// emit unpopulated oneof fields or unpopulated extension fields.
			// The JSON value emitted for unpopulated fields are as follows:
			//  ╔═══════╤════════════════════════════╗
			//  ║ JSON  │ Protobuf field             ║
			//  ╠═══════╪════════════════════════════╣
			//  ║ false │ proto3 boolean fields      ║
			//  ║ 0     │ proto3 numeric fields      ║
			//  ║ ""    │ proto3 string/bytes fields ║
			//  ║ null  │ proto2 scalar fields       ║
			//  ║ null  │ message fields             ║
			//  ║ []    │ list fields                ║
			//  ║ {}    │ map fields                 ║
			//  ╚═══════╧════════════════════════════╝
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			// If DiscardUnknown is set, unknown fields are ignored.
			DiscardUnknown: true,
		},
	})
}

// RegisterForwardHeaders forwards permanent or specified headers from http request to gRPC server.
// Check each `key` within specHeaders, or belong to the list of  permanent request headers maintained by IANA.
// http://www.iana.org/assignments/message-headers/message-headers.xml
func RegisterForwardHeaders(specHeaders map[string]struct{}) gwruntime.ServeMuxOption {
	return gwruntime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		if _, ok := specHeaders[key]; ok {
			return SpecifiedHeaderPrefix + key, ok
		}

		// collect headers in isPermanentHTTPHeader() by default.
		return gwruntime.DefaultHeaderMatcher(key)
	})
}

// RegisterResponseHeaders forwards specified gPRC headers (outgoingHeaders) to http server,
// and inject into response headers.
func RegisterResponseHeaders(outgoingHeaders map[string]struct{}) gwruntime.ServeMuxOption {
	hs := map[string]struct{}{}
	for k := range outgoingHeaders {
		hs[strings.ToLower(k)] = struct{}{}
	}

	return gwruntime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
		if _, ok := hs[key]; ok {
			return key, ok
		}

		return fmt.Sprintf("%s%s", gwruntime.MetadataHeaderPrefix, key), true
	})
}

// enrichMetaData enriches metadata used in gRPC based on http requests.
func enrichMetaData() gwruntime.ServeMuxOption {
	return gwruntime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
		// TODO: inject the identity info based on the auth token.
		route, _ := gwruntime.HTTPPathPattern(ctx)

		return metadata.New(map[string]string{
			HTTPMethod:     r.Method,
			HTTPRequestURI: r.URL.RequestURI(),
			HTTPRealIP:     realip.FromRequest(r),
			HTTPRoute:      route,
			HTTPUserAgent:  r.UserAgent(),
		})
	})
}

// OverrideResponseStatusCode overrides http status code in response header
func OverrideResponseStatusCode() gwruntime.ServeMuxOption {
	return gwruntime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
		code, err := extractStatusCode(ctx)
		if err != nil {
			return err
		}

		if code != 0 {
			// delete the headers to not expose any grpc-metadata in http response
			delete(w.Header(), retHTTPCode)
			w.WriteHeader(code)
		}

		return nil
	})
}

// RegisterHTTPErrorHandler registers customized http handler, overriding http code based on customized status code.
func RegisterHTTPErrorHandler() gwruntime.ServeMuxOption {
	return gwruntime.WithErrorHandler(func(ctx context.Context, mux *gwruntime.ServeMux, m gwruntime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		// override default status code based on native mechanism
		if code, _ := extractStatusCode(ctx); code != 0 {
			err = &gwruntime.HTTPStatusError{HTTPStatus: code, Err: err}
		}

		gwruntime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, err)
	})
}

// InjectHTTPCode injects the customized status code.
func InjectHTTPCode(ctx context.Context, code int) error {
	// TODO: if nedding to customized code, should allow more in the future
	if s := http.StatusText(code); s == "" {
		return ErrInvalidStatusCode
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(HTTPCode, strconv.Itoa(code))); err != nil {
		logger.Ctx(ctx).Error("grpc.SetHeader failed",
			logger.WithError(err),
			logger.WithField("code", code),
		)
		return err
	}

	return nil
}

// InjectHTTPHeader injects customized header with key and value into ctx.
func InjectHTTPHeader(ctx context.Context, key, value string) error {
	if err := grpc.SetHeader(ctx, metadata.Pairs(key, value)); err != nil {
		logger.Ctx(ctx).Error("grpc.SetHeader failed",
			logger.WithError(err),
			logger.WithField("key", key),
			logger.WithField("value", value),
		)
		return err
	}

	return nil
}
