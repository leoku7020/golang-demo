package example

import (
	"context"

	"gorm.io/gorm"

	"demo/internal/models/example"
	exampleUC "demo/internal/usecase/example"
)

type exampleRepo struct {
	db *gorm.DB
}

func NewExampleRepository(db *gorm.DB) exampleUC.ExampleRepository {
	return &exampleRepo{db: db}
}

func (r *exampleRepo) FindMember(ctx context.Context, username string) (*example.Member, error) {
	var member example.Member
	err := r.db.WithContext(ctx).Table("members").Where("username = ?", username).Find(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *exampleRepo) ListItems(ctx context.Context, username, item string) ([]*example.Item, error) {
	var items []*example.Item
	query := r.db.WithContext(ctx).Table("items")

	// TODO 需求不確定
	//if username != "" {
	//	query.Where("name = ?", username)
	//}
	//
	//if item != "" {
	//	query.Where("category = ?", items)
	//}

	err := query.Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
