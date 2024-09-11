package servkit

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"demo/pkg/bootkit"
	"demo/pkg/logger"
)

type RegisterServerFunc func(s *grpc.Server)

// RunGrpcServer starts gRPC server
func RunGrpcServer(
	ctx context.Context, addr string,
	shutdown bootkit.ShutdownFunc,
	registerServers RegisterServerFunc,
	options ...ServOptions,
) error {
	// TODO: only provide tcp now. Need to support unix models socket in the future.
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Ctx(ctx).Error("net.Listen failed",
			logger.WithError(err),
			logger.WithField("addr", addr),
		)
		return err
	}

	o := applyServOptions(options...)
	o.grpcServOpts = append(o.grpcServOpts,
		grpc.ChainUnaryInterceptor(enrichUnaryLogger, LoggingInterceptor),
		grpc.ChainStreamInterceptor(enrichStreamLogger),
	)
	serv := grpc.NewServer(o.grpcServOpts...)
	// register grpc related servers
	registerServers(serv)

	shutdown(func() error {
		// close gRPC server
		logger.Info("Shutting down the gRPC server ...")
		serv.GracefulStop()
		// close listener
		return lis.Close()
	})

	logger.Ctx(ctx).Info(fmt.Sprintf("Starting gRPC server listening at %s", addr))

	return serv.Serve(lis)
}

func enrichUnaryLogger(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := GetMetadata(ctx)
	if !ok {
		return handler(ctx, req)
	}

	if md.Method() == "" {
		return handler(ctx, req)
	}

	var nrTraceID, nrSpanID string
	if txn := newrelic.FromContext(ctx); txn != nil {
		nrMD := txn.GetTraceMetadata()
		nrTraceID = nrMD.TraceID
		nrSpanID = nrMD.SpanID
	}

	if nrTraceID == "" {
		nrTraceID = strings.Replace(uuid.NewString(), "-", "", -1)
	}
	if nrSpanID == "" {
		nrSpanID = strings.Replace(uuid.NewString(), "-", "", -1)
	}

	ctx = logger.ContextWithFields(ctx, logger.Fields{
		"http.method":     md.Method(),
		"http.requestURI": md.RequestURI(),
		"http.route":      md.Route(),
		"http.userAgent":  md.UserAgent(),
		"nr.traceID":      nrTraceID,
		"nr.spanID":       nrSpanID,
	})

	return handler(ctx, req)
}

func enrichStreamLogger(
	srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) error {
	md, ok := GetMetadata(ss.Context())
	if !ok {
		return handler(srv, ss)
	}

	if md.Method() == "" {
		return handler(srv, ss)
	}

	wss := newStreamContextWrapper(ss)

	ctx := wss.Context()
	var nrTraceID, nrSpanID string
	if txn := newrelic.FromContext(ctx); txn != nil {
		nrMD := txn.GetTraceMetadata()
		nrTraceID = nrMD.TraceID
		nrSpanID = nrMD.SpanID
	}

	ctx = logger.ContextWithFields(ctx, logger.Fields{
		"http.method":     md.Method(),
		"http.requestURI": md.RequestURI(),
		"http.route":      md.Route(),
		"http.userAgent":  md.UserAgent(),
		"nr.traceID":      nrTraceID,
		"nr.spanID":       nrSpanID,
	})

	wss.SetContext(ctx)
	return handler(srv, wss)
}

type StreamContextWrapper interface {
	grpc.ServerStream
	SetContext(context.Context)
}

type wrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrapper) Context() context.Context {
	return w.ctx
}

func (w *wrapper) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func newStreamContextWrapper(inner grpc.ServerStream) StreamContextWrapper {
	ctx := inner.Context()
	return &wrapper{
		inner,
		ctx,
	}
}

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	reqMd, _ := metadata.FromIncomingContext(ctx)
	method := reqMd["spec-http-method"]
	apiRoute := reqMd["spec-http-route"]

	resp, err := handler(ctx, req)

	if err != nil {
		logger.Ctx(ctx).Error(fmt.Sprintf("API log: %s %s, Error: %v , Request: %+v", method, apiRoute, err, req))
	} else {
		var ret []byte
		var err error

		if httpBody, ok := resp.(*httpbody.HttpBody); ok {
			ret = httpBody.Data
		} else {
			ret, err = json.Marshal(resp)
		}

		if err != nil {
			logger.Ctx(ctx).Info(fmt.Sprintf("API log: %s %s, Request: %+v, Response: %+v", method, apiRoute, req, resp))
		} else {
			logger.Ctx(ctx).Info(fmt.Sprintf("API log: %s %s, Request: %+v, Response: %+v", method, apiRoute, req, string(ret)))
		}
	}

	return resp, err
}
