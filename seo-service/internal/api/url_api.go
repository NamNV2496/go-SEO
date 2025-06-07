package api

import "github.com/namnv2496/seo/internal/entity"

type CreateUrlRequest struct {
	Url         string                      `query:"url" validate:"required"`
	Name        string                      `query:"name"`
	Tittle      string                      `query:"tittle"`
	Description string                      `query:"description"`
	Template    string                      `query:"template" validate:"required"`
	Prefix      string                      `query:"prefix"`
	Suffix      string                      `query:"suffix"`
	MetaData    []*CreateUrlRequestMetadata `query:"metadata"`
	Domain      string                      `query:"domain"`
	IsActive    bool                        `query:"is_active"`
}

type CreateUrlRequestMetadata struct {
	Keyword string `query:"column:keyword" json:"keyword"`
	Value   string `query:"column:value" json:"value"`
}

type UpdateUrlRequest struct {
	Id          int64                       `param:"id" validate:"required"`
	Url         string                      `query:"url" validate:"required"`
	Name        string                      `query:"name"`
	Tittle      string                      `query:"tittle"`
	Description string                      `query:"description"`
	Template    string                      `query:"template" validate:"required"`
	Prefix      string                      `query:"prefix"`
	Suffix      string                      `query:"suffix"`
	MetaData    []*UpdateUrlRequestMetadata `query:"metadata"`
	Domain      string                      `query:"domain"`
	IsActive    bool                        `query:"is_active"`
}

type UpdateUrlRequestMetadata struct {
	Id      int64  `query:"id" json:"id"`
	UrlId   int64  `query:"url_id" json:"url_id"`
	Keyword string `query:"column:keyword" json:"keyword"`
	Value   string `query:"column:value" json:"value"`
}

type UpdateUrlResponse struct {
	Status string `json:"status"`
}

type GetUrlRequest struct {
	Url string `query:"url" validate:"required"`
}

type GetUrlsRequest struct {
	Page  int `query:"page" validate:"min=1,max=100"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}

type GetUrlsResponse struct {
	Total       int           `json:"total"`
	CurrentPage int           `json:"current_page"`
	Limit       int           `json:"limit"`
	Urls        []*entity.Url `json:"urls"`
}
