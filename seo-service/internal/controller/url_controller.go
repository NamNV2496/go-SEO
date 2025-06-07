package controller

import (
	"log/slog"
	"reflect"
	"strings"

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

	BuildUrl(c echo.Context, req api.BuildUrlRequest) (*api.BuildUrlResponse, error)
	ParseUrl(c echo.Context, req api.ParseUrlRequest) (*api.ParseUrlResponse, error)
	DynamicParamParseByUrl(c echo.Context, req api.DynamicParamRequest) (*api.DynamicParamResponse, error)
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

func (_self *Controller) BuildUrl(c echo.Context, req api.BuildUrlRequest) (*api.BuildUrlResponse, error) {
	ctx := c.Request().Context()
	params := make(map[string]string)
	reqVal := reflect.ValueOf(req)
	for i := 0; i < reqVal.NumField(); i++ {
		field := reqVal.Type().Field(i)
		jsonTag := field.Tag.Get("query")
		if jsonTag == "" {
			continue
		}
		fieldName := strings.Split(jsonTag, ",")[0]
		fieldValue := reqVal.Field(i).Interface()
		if fieldValue != "" {
			params[fieldName] = fieldValue.(string)
		}
	}
	urls, err := _self.urlService.BuildUrl(ctx, req.Kind, params)
	if err != nil {
		return nil, err
	}
	resp := &api.BuildUrlResponse{
		Urls: urls,
	}
	return resp, nil
}

func (_self *Controller) ParseUrl(c echo.Context, req api.ParseUrlRequest) (*api.ParseUrlResponse, error) {
	ctx := c.Request().Context()
	path := strings.Split(req.Url, "/")
	numOfPath := len(path)
	if len(path) == 0 {
		return nil, nil
	}
	urlSeo, err := _self.urlService.ParseUrl(ctx, path[numOfPath-1])
	if err != nil {
		return nil, err
	}
	return &api.ParseUrlResponse{
		Uri:         path[numOfPath-1],
		Path:        "/" + path[numOfPath-1],
		Tittle:      urlSeo.Tittle,
		Description: urlSeo.Description,
	}, nil
}

func (_self *Controller) DynamicParamParseByUrl(c echo.Context, req api.DynamicParamRequest) (*api.DynamicParamResponse, error) {
	ctx := c.Request().Context()
	dynamicParams, err := _self.urlService.DynamicParamParseByUrl(ctx, req.Kind)
	if err != nil {
		return nil, err
	}
	data := make([]*api.DynamicParamData, 0)
	for _, dynamicParam := range dynamicParams {
		data = append(data, &api.DynamicParamData{
			Label: dynamicParam.Name,
			Url:   dynamicParam.Value,
		})
	}
	resp := &api.DynamicParamResponse{
		Data: data,
	}
	return resp, nil
}
