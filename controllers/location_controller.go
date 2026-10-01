package controllers

import (
	"errors"
	"net/http"
	"strings"

	"Event-Explorer/services"

	beego "github.com/beego/beego/v2/server/web"
)

type LocationController struct {
	beego.Controller
	LocationService *services.LocationService
}

func NewLocationController(
	locationService *services.LocationService,
) *LocationController {
	return &LocationController{
		LocationService: locationService,
	}
}


// Autocomplete handles the GET request for location autocomplete suggestions.
func (c *LocationController) Autocomplete() {
	input := c.GetString("input")
	sessionToken := c.GetString("sessionToken")

	suggestions, err := c.LocationService.Autocomplete(
		input,
		sessionToken,
	)

	if err != nil {
		c.CustomAbort(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	c.Data["json"] = map[string]interface{}{
		"suggestions": suggestions,
	}

	c.ServeJSON()
}


// GetPlaceDetails handles the GET request for retrieving place details by place ID.
func (c *LocationController) GetPlaceDetails() {
	placeID := c.Ctx.Input.Param(":placeId")
	sessionToken := c.GetString("sessionToken")

	location, err := c.LocationService.GetPlaceDetails(
		placeID,
		sessionToken,
	)

	if err != nil {
		if errors.Is(
			err,
			services.ErrLocationNotFound,
		) {
			c.CustomAbort(
				http.StatusNotFound,
				err.Error(),
			)
			return
		}

		if strings.TrimSpace(placeID) == "" ||
			sessionToken == "" {
			c.CustomAbort(
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		c.CustomAbort(
			http.StatusBadGateway,
			err.Error(),
		)
		return
	}

	c.Data["json"] = location
	c.ServeJSON()
}