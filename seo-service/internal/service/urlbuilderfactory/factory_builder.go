package urlbuilderfactory

import (
	"context"

	"github.com/namnv2496/seo/internal/entity"
	"gorm.io/gorm"
)

type QueryOption struct {
	Field string
	And   bool // true: and, false: or
}

type IBuilder interface {
	Build(ctx context.Context, request map[string]string) ([]*entity.ShortLink, error)
	BuildRecommend(ctx context.Context, request map[string]string, fields []QueryOption) ([]*entity.ShortLink, error)
}

func BuilderFactory(
	kind string,
	db *gorm.DB,
) (IBuilder, error) {
	switch kind {
	case entity.UrlKindCity:
		return &CityBuilder{
			Db: db,
		}, nil
	case entity.UrlKindProduct:
		return &ProductBuilder{
			Db: db,
		}, nil
	case entity.UrlKindCategory:
		return &CategoryBuilder{
			Db: db,
		}, nil
	case entity.UrlKindBrand:
		return &BrandBuilder{
			Db: db,
		}, nil
	case entity.UrlKindYear:
		return &YearBuilder{
			Db: db,
		}, nil
	default:
		return nil, nil
	}
}
