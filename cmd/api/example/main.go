package main

import (
	"context"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"google.golang.org/grpc"
	"log"
	"os"

	exampleRepo "demo/internal/adapter/repository/mysql/example"
	exampleHlr "demo/internal/router/handler/example"
	exampleUC "demo/internal/usecase/example"
	"demo/pkg/bootkit"
	"demo/pkg/httpkit"
	"demo/pkg/initkit"
	"demo/pkg/logger"
	"demo/pkg/servkit"
	examplePb "demo/proto/example"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// logger
	flushLog := initkit.InitLogger()
	defer flushLog()

	// newrelic.Application
	nrApp := initkit.NewNewRelicApp()

	// mysql
	db := initkit.NewGormDB()

	// migration db in development env
	initkit.ExecMySQLMigration()

	// == middlewares ==
	nrHandler := httpkit.NewNRHttpHandler(nrApp)

	// == business handlers ==
	exampleRepo := exampleRepo.NewExampleRepository(db)
	exampleUC := exampleUC.NewExampleUsecase(exampleRepo)
	exampleHlr := exampleHlr.NewExampleHandler(exampleUC)

	// == run apis ==
	ctx := context.Background()

	grpcAdd := os.Getenv("GRPC_ADDR")
	grpcGWAdd := os.Getenv("GRPC_GW_ADDR")
	bootkit.Register(func(shutdownFn bootkit.ShutdownFunc) error {
		return servkit.RunGrpcServer(ctx, grpcAdd, shutdownFn,
			func(s *grpc.Server) {
				examplePb.RegisterExampleServer(s, exampleHlr)
			},
			servkit.WithGrpcServOptions(
				grpc.ChainUnaryInterceptor(nrgrpc.UnaryServerInterceptor(nrApp)),
				grpc.ChainStreamInterceptor(nrgrpc.StreamServerInterceptor(nrApp)),
			),
		)
	})

	bootkit.Register(func(shutdownFn bootkit.ShutdownFunc) error {
		return servkit.RunGrpcGateway(ctx, grpcGWAdd, grpcAdd, shutdownFn,
			func(ctx context.Context, mux *gwruntime.ServeMux, conn *grpc.ClientConn) error {
				for _, f := range []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{
					examplePb.RegisterExampleHandler,
				} {
					if err := f(ctx, mux, conn); err != nil {
						return err
					}
				}

				return nil
			},
			servkit.WithMiddlewares(nrHandler.Handle),
			servkit.WithGWServMuxOptions(
				servkit.SetFormULREncodeMarhslalerOptions(),
				servkit.SetDefaultMarshalerOptions(),
				servkit.RegisterForwardHeaders(map[string]struct{}{
					"X-Custom-Header": {}, // forward `X-Custom-Header` plus permanent headers from http request into gRPC metadata
				}),
				servkit.OverrideResponseStatusCode(),
				servkit.RegisterHTTPErrorHandler(),
			),
			servkit.WithGrpcDialOptions(
				grpc.WithUnaryInterceptor(nrgrpc.UnaryClientInterceptor),
				grpc.WithStreamInterceptor(nrgrpc.StreamClientInterceptor),
			),
		)
	})

	if err := bootkit.ListenAndServe(ctx); err != nil {
		logger.Ctx(ctx).Fatal("bootkit.ListenAndServe failed", logger.WithError(err))
	}
	// wait until all shutdown callback finished
	bootkit.Wait()
}
