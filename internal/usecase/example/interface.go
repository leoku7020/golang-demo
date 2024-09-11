package example

import (
	"context"

	"demo/internal/models/example"
)

type ExampleRepository interface {
	FindMember(ctx context.Context, username string) (*example.Member, error)
	ListItems(ctx context.Context, username, item string) ([]*example.Item, error)
}

type ExampleUsecase interface {
	Login(ctx context.Context, username, password string) error
	ListItems(ctx context.Context, username, item string) ([]*example.Item, error)
}
