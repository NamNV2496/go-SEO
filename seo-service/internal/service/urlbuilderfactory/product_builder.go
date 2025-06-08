package urlbuilderfactory

import (
	"context"

	"github.com/namnv2496/seo/internal/domain"
	"github.com/namnv2496/seo/internal/entity"
	"github.com/namnv2496/seo/pkg/utils"
	"gorm.io/gorm"
)

type ProductBuilder struct {
	Db *gorm.DB
}

func NewProductBuilder(
	db *gorm.DB,
) *ProductBuilder {
	return &ProductBuilder{
		Db: db,
	}
}

var _ IBuilder = &ProductBuilder{}

func (_self *ProductBuilder) Build(ctx context.Context, request map[string]string) ([]*entity.ShortLink, error) {
	resp := []*entity.ShortLink{}
	err := _self.Db.Model(&domain.ShortLink{}).Where("filter ->> 'product' = ?", request["product"]).Offset(0).Limit(5).Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (_self *ProductBuilder) BuildRecommend(ctx context.Context, request map[string]string, fields []QueryOption) ([]*entity.ShortLink, error) {
	var resp []*entity.ShortLink
	var data []*domain.ShortLink
	// find the same product name
	product := request["product"]
	if product == "" {
		return nil, nil
	}
	tx := _self.Db.Model(&domain.ShortLink{})
	for _, field := range fields {
		if field.And {
			tx = tx.Where("filter->>'"+field.Field+"' =?", request[field.Field])
		} else {
			tx = tx.Or("filter->>'"+field.Field+"' =?", request[field.Field])
		}
	}
	if err := tx.
		Offset(0).
		Limit(5).
		Find(&data).Error; err != nil {
		return nil, err
	}
	utils.Copy(&resp, data)
	return resp, nil
}
