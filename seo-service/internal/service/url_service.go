package service

import (
	"context"

	"github.com/namnv2496/seo/internal/domain"
	"github.com/namnv2496/seo/internal/entity"
	"github.com/namnv2496/seo/internal/repository"
	"github.com/namnv2496/seo/pkg/utils"
	"gorm.io/gorm"
)

type IUrlService interface {
	ParseUrl(ctx context.Context, url string) (*entity.Url, error)
	BuildUrl(ctx context.Context, request map[string]string) (string, error)

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
}

func NewUrlService(
	db repository.IDatabase,
	urlRepo repository.IUrlRepository,
	urlMetadataRepo repository.IUrlMetadataRepository,
) *UrlService {
	return &UrlService{
		db:              db,
		urlRepo:         urlRepo,
		urlMetadataRepo: urlMetadataRepo,
	}
}

var _ IUrlService = &UrlService{}

func (_self *UrlService) ParseUrl(ctx context.Context, url string) (*entity.Url, error) {
	return nil, nil
}

func (_self *UrlService) BuildUrl(ctx context.Context, request map[string]string) (string, error) {
	return "", nil
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
