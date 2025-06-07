package controller

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/namnv2496/seo/internal/api"
	"github.com/namnv2496/seo/internal/entity"
	"github.com/namnv2496/seo/internal/service"
	"github.com/namnv2496/seo/pkg/utils"
)

type IController interface {
	CreateNewUrl(c echo.Context, req api.CreateUrlRequest) error
	UpdateUrl(c echo.Context, req api.UpdateUrlRequest) (*api.UpdateUrlResponse, error)
	GetUrl(c echo.Context, req api.GetUrlRequest) (*entity.Url, error)
	GetUrls(c echo.Context, req api.GetUrlsRequest) (*api.GetUrlsResponse, error)
}

type Controller struct {
	urlService service.IUrlService
}

func NewUrlController(
	urlService service.IUrlService,
) *Controller {
	return &Controller{
		urlService: urlService,
	}
}

var _ IController = &Controller{}

func (_self *Controller) CreateNewUrl(c echo.Context, req api.CreateUrlRequest) error {
	slog.Info("CreateNewUrl", "req", req)
	ctx := c.Request().Context()
	var request entity.Url
	utils.Copy(&request, req)
	err := _self.urlService.CreateUrl(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (_self *Controller) UpdateUrl(c echo.Context, req api.UpdateUrlRequest) (*api.UpdateUrlResponse, error) {
	ctx := c.Request().Context()
	var request entity.Url
	utils.Copy(&request, req)
	err := _self.urlService.UpdateUrl(ctx, request)
	if err != nil {
		return nil, err
	}
	return &api.UpdateUrlResponse{
		Status: "success",
	}, nil
}

func (_self *Controller) GetUrl(c echo.Context, req api.GetUrlRequest) (*entity.Url, error) {
	ctx := c.Request().Context()
	url, err := _self.urlService.GetUrl(ctx, req.Url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (_self *Controller) GetUrls(c echo.Context, req api.GetUrlsRequest) (*api.GetUrlsResponse, error) {
	ctx := c.Request().Context()
	urls, err := _self.urlService.GetUrls(ctx, (req.Page-1)*req.Limit, req.Limit)
	if err != nil {
		return nil, err
	}
	resp := &api.GetUrlsResponse{
		Total:       len(urls),
		CurrentPage: req.Page,
		Limit:       req.Limit,
		Urls:        urls,
	}
	return resp, nil
}
