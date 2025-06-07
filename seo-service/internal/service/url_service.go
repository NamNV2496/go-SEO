package service

import (
	"bytes"
	"context"
	"html/template"
	"regexp"
	"strings"

	"github.com/namnv2496/seo/configs"
	"github.com/namnv2496/seo/internal/domain"
	"github.com/namnv2496/seo/internal/entity"
	"github.com/namnv2496/seo/internal/repository"
	"github.com/namnv2496/seo/pkg/utils"
	"gorm.io/gorm"
)

type IUrlService interface {
	ParseUrl(ctx context.Context, url string) (*entity.Url, error)
	BuildUrl(ctx context.Context, kind string, request map[string]string) ([]string, error)
	DynamicParamParseByUrl(ctx context.Context, url string) ([]*entity.DynamicParam, error)

	CreateUrl(ctx context.Context, url entity.Url) error
	UpdateUrl(ctx context.Context, url entity.Url) error
	DeleteUrl(ctx context.Context, url entity.Url) error
	GetUrl(ctx context.Context, url string) (*entity.Url, error)
	GetUrls(ctx context.Context, offset, limit int) ([]*entity.Url, error)
}

type UrlService struct {
	db              repository.IDatabase
	urlRepo         repository.IUrlRepository
	urlMetadataRepo repository.IUrlMetadataRepository
	buildUrlTool    string
}

func NewUrlService(
	conf *configs.Config,
	db repository.IDatabase,
	urlRepo repository.IUrlRepository,
	urlMetadataRepo repository.IUrlMetadataRepository,
) *UrlService {
	return &UrlService{
		db:              db,
		urlRepo:         urlRepo,
		urlMetadataRepo: urlMetadataRepo,
		buildUrlTool:    conf.BuildUrlTool,
	}
}

var _ IUrlService = &UrlService{}

func (_self *UrlService) ParseUrl(ctx context.Context, url string) (*entity.Url, error) {
	urlData, err := _self.urlRepo.GetUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	if urlData == nil {
		return nil, nil
	}
	metadata, err := _self.urlMetadataRepo.GetUrlMetadata(ctx, urlData.Id)
	if err != nil {
		return nil, err
	}
	params := make(map[string]string, 0)
	resp := &entity.Url{}
	utils.Copy(&resp, urlData)
	if metadata != nil {
		var metadataEntity []*entity.UrlMetadata
		for _, meta := range metadata {
			var elem *entity.UrlMetadata
			utils.Copy(&elem, meta)
			metadataEntity = append(metadataEntity, elem)
			params[meta.Keyword] = meta.Value
		}
		resp.Metadata = metadataEntity
	}
	resp.Tittle, _ = utils.BuildByTemplate(ctx, "build-tittle", urlData.Template, params)
	resp.Description, _ = utils.BuildByTemplate(ctx, "build-description", urlData.Template, params)
	return resp, nil
}

func (_self *UrlService) BuildUrl(ctx context.Context, kind string, request map[string]string) ([]string, error) {
	// TBU: ==================== Intergrate AI to build ====================
	switch _self.buildUrlTool {
	case "ai":
		return buildUrlByAI(ctx, kind, request)
	case "template":
		return _self.buildUrlByTemplate(ctx, kind, request)
	case "regex":
		return buildUrlByRegex(ctx, kind, request)
	default:
		return []string{"not-found"}, nil
	}
}

func buildUrlByAI(ctx context.Context, kind string, request map[string]string) ([]string, error) {
	return []string{}, nil
}

func (_self *UrlService) buildUrlByTemplate(ctx context.Context, kind string, request map[string]string) ([]string, error) {
	// Construct the URL template
	urlData, err := _self.urlRepo.GetUrl(ctx, kind)
	if err != nil {
		return []string{}, err
	}
	if urlData == nil {
		return []string{}, nil
	}
	urlTemplate := template.New("url-template")
	urlText := urlData.Template
	if urlData.Prefix != "" {
		urlText = urlData.Prefix + urlData.Template
	}
	if urlData.Suffix != "" {
		urlText = urlText + "-" + urlData.Suffix
	}
	urlTemplate, err = urlTemplate.Parse(urlText)
	if err != nil {
		return []string{}, err
	}
	result := new(bytes.Buffer)
	request["kind"] = kind
	if err := urlTemplate.Execute(result, request); err != nil {
		return []string{}, err
	}
	return []string{result.String()}, nil
}

func buildUrlByRegex(ctx context.Context, kind string, request map[string]string) ([]string, error) {
	var templates = []string{
		"{kind}-{location}-{category}-{product}-{brand}",
		"{kind}-{category}-{product}",
		"{kind}-{location}-{product}",
		"{kind}-{category}-{brand}-{product}",
		"{kind}-{category}-{year}-{month}",
		"{kind}-{product}",
	}
	request["kind"] = kind
	var urls []string
	for _, tpl := range templates {
		requiredFields := extractor(tpl)
		skip := false
		for _, field := range requiredFields {
			if val, ok := request[field]; !ok || val == "" {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		url := tpl
		for key, val := range request {
			url = strings.ReplaceAll(url, "{"+key+"}", val)
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func extractor(template string) []string {
	re := regexp.MustCompile(`\{(\w+)\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	var fields []string
	for _, match := range matches {
		fields = append(fields, match[1])
	}
	return fields
}

func (_self *UrlService) DynamicParamParseByUrl(ctx context.Context, url string) ([]*entity.DynamicParam, error) {
	return nil, nil
}

func (_self *UrlService) CreateUrl(ctx context.Context, url entity.Url) error {
	var request domain.Url
	utils.Copy(&request, url)
	var newUrlId int64
	err := _self.db.RunWithTransaction(ctx,
		func(ctx context.Context, tx *gorm.DB) error {
			urlId, err := _self.urlRepo.CreateUrl(ctx, tx, request)
			if err != nil {
				return err
			}
			newUrlId = urlId
			return nil
		},
		func(ctx context.Context, tx *gorm.DB) error {
			var metadata []*domain.UrlMetadata
			for _, meta := range url.Metadata {
				var elem *domain.UrlMetadata
				utils.Copy(&elem, meta)
				elem.UrlId = newUrlId
				metadata = append(metadata, elem)
			}
			if err := _self.urlMetadataRepo.CreateUrlMetadata(ctx, tx, metadata); err != nil {
				return err
			}
			return nil
		})
	return err
}

func (_self *UrlService) UpdateUrl(ctx context.Context, url entity.Url) error {
	err := _self.db.RunWithTransaction(ctx,
		func(ctx context.Context, tx *gorm.DB) error {
			var request domain.Url
			utils.Copy(&request, url)
			if err := _self.urlRepo.UpdateUrl(ctx, tx, request); err != nil {
				return err
			}
			return nil
		},
		func(ctx context.Context, tx *gorm.DB) error {
			var metadata []*domain.UrlMetadata
			for _, meta := range url.Metadata {
				var elem *domain.UrlMetadata
				utils.Copy(&elem, meta)
				metadata = append(metadata, elem)
			}
			if err := _self.urlMetadataRepo.UpdateUrlMetadata(ctx, tx, metadata); err != nil {
				return err
			}
			return nil
		})
	return err
}

func (_self *UrlService) DeleteUrl(ctx context.Context, url entity.Url) error {
	err := _self.db.RunWithTransaction(ctx,
		func(ctx context.Context, tx *gorm.DB) error {
			var request domain.Url
			utils.Copy(&request, url)
			if err := _self.urlRepo.DeleteUrl(ctx, tx, request.Url); err != nil {
				return err
			}
			return nil
		},
		func(ctx context.Context, tx *gorm.DB) error {
			if err := _self.urlMetadataRepo.DeleteUrlMetadataById(ctx, tx, url.Id); err != nil {
				return err
			}
			return nil
		})
	return err
}

func (_self *UrlService) GetUrl(ctx context.Context, url string) (*entity.Url, error) {
	urlData, err := _self.urlRepo.GetUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	if urlData == nil {
		return nil, nil
	}
	metadata, err := _self.urlMetadataRepo.GetUrlMetadata(ctx, urlData.Id)
	if err != nil {
		return nil, err
	}
	resp := &entity.Url{}
	utils.Copy(&resp, urlData)
	if metadata != nil {
		var metadataEntity []*entity.UrlMetadata
		for _, meta := range metadata {
			var elem *entity.UrlMetadata
			utils.Copy(&elem, meta)
			metadataEntity = append(metadataEntity, elem)
		}
		resp.Metadata = metadataEntity
	}
	return resp, nil
}

func (_self *UrlService) GetUrls(ctx context.Context, offset, limit int) ([]*entity.Url, error) {
	urlDatas, err := _self.urlRepo.GetUrls(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, urlData := range urlDatas {
		ids = append(ids, urlData.Id)
	}
	metadata, err := _self.urlMetadataRepo.GetUrlMetadatas(ctx, ids)
	if err != nil {
		return nil, err
	}
	var resp []*entity.Url
	for _, urlData := range urlDatas {
		var metadataEntity []*entity.UrlMetadata
		for _, meta := range metadata {
			var elem *entity.UrlMetadata
			utils.Copy(&elem, meta)
			metadataEntity = append(metadataEntity, elem)
		}
		var urlElem *entity.Url
		utils.Copy(&resp, urlData)
		urlElem = &entity.Url{}
		utils.Copy(&urlElem, urlData)
		urlElem.Metadata = metadataEntity
		resp = append(resp, urlElem)
	}
	return resp, nil
}
