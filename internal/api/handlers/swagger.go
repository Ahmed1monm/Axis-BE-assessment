package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// ServeUI handles serving the Swagger UI
func (h *SwaggerHandler) ServeUI(c echo.Context) error {
	return echoSwagger.WrapHandler(c)
}

// ServeSpec handles serving the swagger.yaml specification file
func (h *SwaggerHandler) ServeSpec(c echo.Context) error {
	return c.File(filepath.Join("docs", "swagger.yaml"))
}

// RedirectToUI handles redirecting /swagger to /swagger/index.html
func (h *SwaggerHandler) RedirectToUI(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
}
