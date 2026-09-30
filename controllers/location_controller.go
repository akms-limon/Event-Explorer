package controllers

import (
	"net/http"

	"Event-Explorer/services"

	beego "github.com/beego/beego/v2/server/web"
)

type LocationController struct {
	beego.Controller
	locationService *services.LocationService
}

func NewLocationController(
	locationService *services.LocationService,
) *LocationController {
	return &LocationController{
		locationService: locationService,
	}
}

func (c *LocationController) Autocomplete() {
	input := c.GetString("input")
	sessionToken := c.GetString("sessionToken")

	suggestions, err := c.locationService.Autocomplete(
		c.Ctx.Request.Context(),
		input,
		sessionToken,
	)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"suggestions": suggestions,
	}

	c.ServeJSON()
}

func (c *LocationController) GetPlaceDetails() {
	placeID := c.Ctx.Input.Param(":placeId")
	sessionToken := c.GetString("sessionToken")

	location, err := c.locationService.GetPlaceDetails(
		c.Ctx.Request.Context(),
		placeID,
		sessionToken,
	)
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, err.Error())
		return
	}

	c.Data["json"] = location

	c.ServeJSON()
}
