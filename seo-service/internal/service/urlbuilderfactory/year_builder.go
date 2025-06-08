package urlbuilderfactory

import (
	"context"
	"fmt"
	"strconv"

	"github.com/namnv2496/seo/internal/domain"
	"github.com/namnv2496/seo/internal/entity"
	"gorm.io/gorm"
)

type YearBuilder struct {
	Db *gorm.DB
}

func NewYearBuilder(
	db *gorm.DB,
) *YearBuilder {
	return &YearBuilder{
		Db: db,
	}
}

var _ IBuilder = &YearBuilder{}

func (_self *YearBuilder) Build(ctx context.Context, request map[string]string) ([]*entity.ShortLink, error) {
	resp := []*entity.ShortLink{}
	err := _self.Db.Model(&domain.ShortLink{}).Where("filter -> year = ?", request["year"]).Offset(0).Limit(5).Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (_self *YearBuilder) BuildRecommend(ctx context.Context, request map[string]string, fields []QueryOption) ([]*entity.ShortLink, error) {
	var resp []*entity.ShortLink
	var nextYears []*entity.ShortLink
	var PrevYears []*entity.ShortLink
	yearText := request["year"]
	if yearText == "" {
		return nil, nil
	}
	var err error
	var year int64
	year, err = strconv.ParseInt(yearText, 10, 64)
	if err != nil {
		return nil, err
	}
	// find next years
	err = _self.Db.Model(&domain.ShortLink{}).Where("filter ->> 'year' = ?", fmt.Sprintf("%d", year+1)).Offset(0).Limit(5).Find(&nextYears).Error
	if err != nil {
		return nil, err
	}
	// find prev years
	err = _self.Db.Model(&domain.ShortLink{}).Where("filter ->> 'year' = ?", fmt.Sprintf("%d", year-1)).Offset(0).Limit(5).Find(&PrevYears).Error
	if err != nil {
		return nil, err
	}

	resp = append(resp, nextYears...)
	resp = append(resp, PrevYears...)
	return resp, nil
}
