package server

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	handlers "github.com/subannn/urlshorter/handlers"
	"golang.org/x/net/context"
)

var e *echo.Echo

func RunServer() {
	// Echo instance
	e = echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/*", handlers.RedirectToLongURL)
	e.POST("/cutURL", handlers.CutLongURL)

	// Start server

	if err := e.Start(":1323"); err != nil {
		log.Println(err)
	}
	
}

func ShutDownServer(ctxToShutdown context.Context) {
	e.Shutdown(ctxToShutdown)
}
