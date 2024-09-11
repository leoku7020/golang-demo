package example

import (
	"context"
	"crypto/md5"
	"demo/internal/models/example"
	"encoding/hex"
	"errors"
	"strings"
)

func NewExampleUsecase(repo ExampleRepository) ExampleUsecase {
	return &impl{
		repo: repo,
	}
}

type impl struct {
	repo ExampleRepository
}

func (im *impl) Login(ctx context.Context, username, password string) error {
	member, err := im.repo.FindMember(ctx, username)
	if err != nil {
		return err
	}
	hash := md5.New()
	hash.Write([]byte(password))
	hashInBytes := hash.Sum(nil)
	hashInHex := hex.EncodeToString(hashInBytes)
	if strings.ToLower(member.Password) != hashInHex {
		return errors.New("wrong password")
	}
	return nil
}

func (im *impl) ListItems(ctx context.Context, username, item string) ([]*example.Item, error) {
	items, err := im.repo.ListItems(ctx, username, item)
	if err != nil {
		return nil, err
	}
	return items, nil
}
