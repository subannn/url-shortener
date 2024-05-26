package server

import (
  "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
  handlers "github.com/subannn/urlshorter/handlers"
)

func RunServer() {
  // Echo instance
  e := echo.New()
  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // Routes
  e.GET("/*", handlers.RedirectToLongURL)
  e.POST("/cutURL", handlers.CutLongURL)

  // Start server
  e.Logger.Fatal(e.Start(":1323"))
}