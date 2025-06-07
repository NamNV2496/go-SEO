package service

import (
	"context"

	"github.com/namnv2496/seo/internal/entity"
	"github.com/namnv2496/seo/internal/repository"
)

func buildDynamic(kind string, request map[string]string) (map[string]string, error) {
	var resp map[string]string
	switch kind {
	case entity.UrlKindLocation:
		resp = map[string]string{
			"location_id": request["location_id"],
		}
	case entity.UrlKindProduct:
		resp = map[string]string{
			"product_id": request["product_id"],
		}
	case entity.UrlKindCategory:
		resp = map[string]string{
			"category_id": request["category_id"],
		}
	case entity.UrlKindBrand:
		resp = map[string]string{
			"brand_id": request["brand_id"],
		}
	case entity.UrlKindYear:
		resp = map[string]string{
			"year": request["year"],
		}
	case entity.UrlKindMonth:
		resp = map[string]string{
			"month": request["month"],
		}
	}

	return resp, nil
}

type IBuilder interface {
	Build(ctx context.Context, kind string, request map[string]string) (string, error)
}

func BuilderFactory(
	kind string,
	repository repository.IUrlRepository,
) (IBuilder, error) {
	// switch kind {
	// case entity.UrlKindLocation:
	// 	return &LocationBuilder{
	// 		repository: repository,
	// 	}, nil
	// case entity.UrlKindProduct:
	// 	return &ProductBuilder{
	// 		repository: repository,
	// 	}, nil
	// case entity.UrlKindCategory:
	// 	return &CategoryBuilder{
	// 		repository: repository,
	// 	}, nil
	// case entity.UrlKindBrand:
	// 	return &BrandBuilder{
	// 		repository: repository,
	// 	}, nil
	// case entity.UrlKindYear:
	// 	return &YearBuilder{
	// 		repository: repository,
	// 	}, nil
	// case entity.UrlKindMonth:
	// 	return &MonthBuilder{
	// 		repository: repository,
	// 	}, nil
	// default:
	// 	return nil, nil
	// }
	return nil, nil
}
