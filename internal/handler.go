package internal

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func parseLinkResult(data []LinkResponse, err error) *LinkResponse {
	if err == nil && len(data) > 0 {
		var id = data[0].ID
		_, consumeError := ConsumeDeeplink(id)

		if consumeError != nil {
			fmt.Println("Error when trying to consume id " + strconv.Itoa(id) + " " + consumeError.Error())
		}

		if len(data) > 1 {
			// More than one found... this is a problem...lets log it
			fmt.Println("Found more than one link for id: " + strconv.Itoa(id))
		}

		return &data[0]
	}

	return nil
}

func GetLinkFromId(c echo.Context) error {
	// id here is a string but we don't know if this is a fingerprint or a memberIdHash. We will check against both.
	// TODO: In the future we will want to only get the memberId from the JWT token.
	paramId := c.Param("paramId")

	// Prioritize data retrieved by memberIdHash
	data, err := GetLinksByMemberIdHash(paramId)
	result := parseLinkResult(data, err)
	if err == nil && result != nil {
		return c.JSON(http.StatusOK, result)
	}

	// Fall back on data retrieved by fingerprint
	data, err = GetLinksByFingerprint(paramId)
	result = parseLinkResult(data, err)
	if err == nil && result != nil {
		return c.JSON(http.StatusOK, result)
	}

	if err != nil {
		// TODO handle error
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

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

// StoreDeeplink /**
/**
 * Store a deeplink that might come in from a targeted email campaign.
 */
func StoreDeeplink(c echo.Context) error {
	// TODO: Hacky way around issue
	c.Request().Header.Set("Content-Type", "application/json")
	ls := new(LinkStoreRequest)
	if err := c.Bind(ls); err != nil {
		return c.JSON(http.StatusBadRequest, "Missing required fields")
	}

	err := SaveLink(ls)
	if err != nil && err.Error() == "missing required property" {
		return c.JSON(http.StatusBadRequest, "Missing required fields")
	} else if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Stored link for member: " + ls.MemberIdHash + " with fingerprint: " + ls.Fingerprint)
	return c.JSON(http.StatusOK, "Link stored")
}
