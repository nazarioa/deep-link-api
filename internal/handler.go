package internal

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetLinkByFingerprint(c echo.Context) error {
	fingerPrint := c.Param("fingerprint")
	data, err := GetLinksByFingerprint(fingerPrint)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if len(data) == 1 {
		return c.JSON(http.StatusOK, data[0])
	} else if len(data) > 1 {
		fmt.Println(len(data))
		fmt.Println(data)
		return c.JSON(http.StatusConflict, "Multiple links found for the same fingerprint")
	}
	return c.JSON(http.StatusNoContent, data)
}
