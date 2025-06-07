package controller

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namnv2496/seo/configs"
	"github.com/namnv2496/seo/pkg"
)

func Start(
	urlController IController,
) (*echo.Echo, error) {
	conf := configs.LoadConfig()
	e := newEchoServer()
	publicGroup := e.Group("/api/v1/public")
	publicGroup.POST("/url", wrapReponse(urlController.CreateNewUrl))
	publicGroup.PUT("/url/:id", wrapReponse(urlController.UpdateUrl))
	publicGroup.GET("/url", wrapReponse(urlController.GetUrl))
	publicGroup.GET("/urls", wrapReponse(urlController.GetUrls))

	if err := e.Start(fmt.Sprintf(":%s", conf.AppPort)); err != nil {
		e.Logger.Fatal(err)
	}
	slog.Info("Server is running on port: ", "port", conf.AppPort)
	return e, nil
}

func newEchoServer() *echo.Echo {
	e := echo.New()
	e.Validator = pkg.NewValidator()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	return e
}

func wrapReponse(function any) echo.HandlerFunc {
	ftype := reflect.TypeOf(function)
	fval := reflect.ValueOf(function)
	if ftype.NumIn() != 2 {
		panic("function must have 2 parameters")
	}
	if fval.Kind() != reflect.Func {
		panic("function must be a function")
	}
	runtime.FuncForPC(fval.Pointer()).Name()
	errorIndex := ftype.NumOut() - 1

	return func(c echo.Context) error {
		// execute function
		req := reflect.New(ftype.In(1))
		if err := c.Bind(req.Interface()); err != nil {
			return err
		}
		err := c.Validate(req.Interface())
		if err != nil {
			return err
		}

		res := fval.Call([]reflect.Value{
			reflect.ValueOf(c),
			req.Elem(),
		})
		if !res[errorIndex].IsNil() {
			return res[errorIndex].Interface().(error)
		}
		resp := c.Response()
		output := res[0].Interface()
		return c.JSON(resp.Status, output)
	}
}
