package example

import (
	"context"
	"net/http"

	"demo/internal/router/handler"
	uc "demo/internal/usecase/example"
	"demo/pkg/logger"
	pb "demo/proto/example"
)

type exampleHandler struct {
	*pb.UnimplementedExampleServer

	usecase uc.ExampleUsecase
}

func NewExampleHandler(usecase uc.ExampleUsecase) pb.ExampleServer {
	return &exampleHandler{usecase: usecase}
}

func (h *exampleHandler) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	if err := req.Validate(); err != nil {
		logger.Ctx(ctx).Error("Validate failed", logger.WithError(err))
		return handler.AbortWithError(ctx, err, &pb.LoginResp{
			Status:      "failed",
			Description: "bad request",
		},
			handler.WithHttpStatus(http.StatusBadRequest),
		)
	}

	if err := h.usecase.Login(ctx, req.Username, req.Password); err != nil {
		logger.Ctx(ctx).Error("usecase.Login failed", logger.WithError(err))
		return handler.AbortWithError(ctx, err, &pb.LoginResp{
			Status:      "failed",
			Description: "wrong password",
		}, handler.WithHttpStatus(http.StatusBadRequest))
	}

	return handler.RenderResponse(ctx, &pb.LoginResp{
		Status:      "success",
		Description: "ok",
	},
		handler.WithHttpStatus(http.StatusOK),
	)
}

func (h *exampleHandler) ListItems(ctx context.Context, req *pb.ListItemsReq) (*pb.ListItemsResp, error) {
	if err := req.Validate(); err != nil {
		logger.Ctx(ctx).Error("Validate failed", logger.WithError(err))
		return handler.AbortWithError(ctx, err, &pb.ListItemsResp{
			Status: "failed",
			Item:   nil,
		},
			handler.WithHttpStatus(http.StatusBadRequest),
		)
	}

	items, err := h.usecase.ListItems(ctx, req.Username, req.Item)
	if err != nil {
		logger.Ctx(ctx).Error("usecase.ListRecords failed", logger.WithError(err))
		return handler.AbortWithError(ctx, err, &pb.ListItemsResp{})
	}

	rs := make([]*pb.ItemData, len(items))
	for i, r := range items {
		rs[i] = &pb.ItemData{
			ItemId:   r.Id,
			ItemName: r.Name,
			Category: r.Category,
		}
	}

	return handler.RenderResponse(ctx, &pb.ListItemsResp{Status: "success", Item: rs})
}
