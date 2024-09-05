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

// GetLinkByMemberIdHash /**
/**
 * Given a memberIdHash, return the link that was stored for that member.
 */
func GetLinkByMemberIdHash(c echo.Context) error {
	memberIdHash := c.Param("memberIdHash")
	// TODO confirm that the memberIdHash is for the currently authenticated member

	data, err := GetLinksByMemberIdHash(memberIdHash)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if len(data) == 1 {
		return c.JSON(http.StatusOK, data[0])
	} else if len(data) > 1 {
		fmt.Println(len(data))
		fmt.Println(data)
		return c.JSON(http.StatusConflict, "Multiple links found for the same memberId")
	}
	return c.JSON(http.StatusNoContent, data)
}

func StoreDeeplink(c echo.Context) error {
	l := new(Link)

	err := SaveLink(l)
	if err != nil && err.Error() == "missing required property" {
		return c.JSON(http.StatusBadRequest, "Missing required fields")
	} else if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Stored link for member: " + l.MemberIdHash + " with fingerprint: " + l.Fingerprint)
	return c.JSON(http.StatusOK, "Link stored")
}
