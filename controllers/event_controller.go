package controllers

import (
	"errors"
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

	musicEvents, sportsEvents, err := c.EventService.GetEvents(
		city,
		countryCode,
	)

	if err != nil {
		c.CustomAbort(http.StatusBadGateway, err.Error())
		return
	}

	c.Data["City"] = city
	c.Data["CountryCode"] = countryCode
	c.Data["MusicEvents"] = musicEvents
	c.Data["SportsEvents"] = sportsEvents
	c.TplName = "listing.tpl"
}

func (c *EventController) Details() {
	eventID := c.Ctx.Input.Param(":eventId")

	event, err := c.EventService.GetEvent(eventID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidEventID) {
			c.CustomAbort(http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, services.ErrEventNotFound) {
			c.CustomAbort(http.StatusNotFound, err.Error())
			return
		}

		c.CustomAbort(http.StatusBadGateway, err.Error())
		return
	}

	c.Data["Event"] = event
	c.TplName = "details.tpl"
}

func (c *EventController) Redirect() {
	eventID := c.Ctx.Input.Param(":eventId")

	ticketURL, err := c.EventService.GetTicketURL(eventID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidEventID) {
			c.CustomAbort(http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, services.ErrEventNotFound) {
			c.CustomAbort(http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, services.ErrTicketURLMissing) ||
			errors.Is(err, services.ErrUnsafeTicketURL) {
			c.CustomAbort(http.StatusBadRequest, err.Error())
			return
		}

		c.CustomAbort(http.StatusBadGateway, err.Error())
		return
	}

	c.Controller.Redirect(
		ticketURL,
		http.StatusFound,
	)
}
