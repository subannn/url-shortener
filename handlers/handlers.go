package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	models "github.com/subannn/urlshorter/models"
	redis "github.com/subannn/urlshorter/redis"
)

// Handler
func CutLongURL(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		panic(err)
	}

	defer c.Request().Body.Close()
	var URL models.RequestLongURL
	err = json.Unmarshal(body, &URL)
	if err != nil {
		panic(err)
	}

	if(len(URL.LongURL) > 2000) {
		panic("Input URL length exceeds 2000")
	}

	shortURL := redis.CutAndSaveURL(URL.LongURL)

	return c.String(http.StatusOK, shortURL)
}

func RedirectToLongURL(c echo.Context) error {
	shortURL := c.Request().URL.Path[1:] // remove '/' at the begining

	longURL := redis.GetLongURL(shortURL)

	return c.Redirect(http.StatusPermanentRedirect, longURL)
}
