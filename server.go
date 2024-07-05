package main

import (
	"deep-link-api/internal"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	// TODO: Get env variable and pass to InitDb
	err := internal.InitDb("")

	// TODO: defer close on panic?

	// This is a panic because the server should not start if the database is not connected
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/f/:fingerprint", internal.GetLinkByFingerprint)
	e.Logger.Fatal(e.Start(":1323"))
}
