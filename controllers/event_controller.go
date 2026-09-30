package controllers

import (
	"net/http"

	"Event-Explorer/services"

	beego "github.com/beego/beego/v2/server/web"
)

type EventController struct {
	beego.Controller
	EventService *services.EventService
}

func NewEventController(
	eventService *services.EventService,
) *EventController {
	return &EventController{
		EventService: eventService,
	}
}

func (c *EventController) List() {
	city := c.GetString("city")
	countryCode := c.GetString("countryCode")

	musicEvents, sportsEvents, err :=
		c.EventService.GetEvents(
			city,
			countryCode,
		)

	if err != nil {
		c.CustomAbort(
			http.StatusBadGateway,
			err.Error(),
		)
		return
	}

	c.Data["City"] = city
	c.Data["CountryCode"] = countryCode
	c.Data["MusicEvents"] = musicEvents
	c.Data["SportsEvents"] = sportsEvents

	c.TplName = "listing.tpl"
}