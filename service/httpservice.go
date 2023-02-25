package service

import (
	"net/http"
	"silvernote/factory"
	"silvernote/service/handler"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HttpServer struct {
	Fac            *factory.Factory
	servicehandler handler.ServiceHandler
}

func (_self *HttpServer) StartService() {
	_self.servicehandler = handler.ServiceHandler{Fac: _self.Fac}

	e := echo.New()

	store := sessions.NewCookieStore([]byte("WEB"))

	store.Options = &sessions.Options{
		Path:     "/",
		Secure:   true,
		MaxAge:   60,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(store))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	e.POST("/service", _self.servicehandler.RequestHandle)

	e.Logger.Fatal(e.Start(":1000"))
}

func (_self *HttpServer) Monitoring(c echo.Context) error {
	c.HTML(200, "")
	return nil
}
