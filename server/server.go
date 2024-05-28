package server

import (
	"log"
	"os"

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

	port := ":" + os.Getenv("SERVER_PORT")
	if err := e.Start(port); err != nil {
		log.Println(err)
	}
	
}

func ShutDownServer(ctxToShutdown context.Context) {
	err := e.Shutdown(ctxToShutdown)
	if err != nil {
		log.Println(err)
	}
}
