package controllers

import (
	"errors"
	"net/http"
	"strings"

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

	musicEvents,
		sportsEvents,
		musicCacheHit,
		sportsCacheHit,
		err := c.EventService.GetEvents(
		city,
		countryCode,
	)

	if err != nil {
		if errors.Is(
			err,
			services.ErrInvalidEventSearchInput,
		) {
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

	c.Data["City"] = city
	c.Data["CountryCode"] = countryCode

	c.Data["MusicEvents"] = musicEvents
	c.Data["SportsEvents"] = sportsEvents

	c.Data["MusicCacheHit"] = musicCacheHit
	c.Data["SportsCacheHit"] = sportsCacheHit

	c.TplName = "listing.tpl"
}

func (c *EventController) Details() {
	eventID := c.Ctx.Input.Param(":eventId")

	event, err := c.EventService.GetEvent(eventID)

	if err != nil {
		if errors.Is(
			err,
			services.ErrInvalidEventID,
		) {
			c.CustomAbort(
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		if errors.Is(
			err,
			services.ErrEventNotFound,
		) {
			c.CustomAbort(
				http.StatusNotFound,
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

	c.Data["Event"] = event
	c.TplName = "details.tpl"
}

func (c *EventController) Redirect() {
	eventID := c.Ctx.Input.Param(":eventId")

	ticketURL, err := c.EventService.GetTicketURL(eventID)

	if err != nil {
		if errors.Is(
			err,
			services.ErrInvalidEventID,
		) {
			c.CustomAbort(
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		if errors.Is(
			err,
			services.ErrEventNotFound,
		) {
			c.CustomAbort(
				http.StatusNotFound,
				err.Error(),
			)
			return
		}

		if errors.Is(
			err,
			services.ErrTicketURLMissing,
		) {
			c.CustomAbort(
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		if errors.Is(
			err,
			services.ErrUnsafeTicketURL,
		) {
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

	c.Controller.Redirect(
		ticketURL,
		http.StatusFound,
	)
}

func (c *EventController) SearchCachedCities() {
	search := strings.TrimSpace(
		c.GetString("search"),
	)

	if len([]byte(search)) < 3 {
		c.Data["json"] = map[string]interface{}{
			"locations": []services.CacheLocation{},
		}

		c.ServeJSON()
		return
	}

	locations := c.EventService.GetCachedLocations(search)

	c.Data["json"] = map[string]interface{}{
		"locations": locations,
	}

	c.ServeJSON()
}

func (c *EventController) InvalidateCache() {
	c.EventService.InvalidateCache()

	c.Data["json"] = map[string]interface{}{
		"message": "all event cache cleared",
	}

	c.ServeJSON()
}

func (c *EventController) InvalidateCacheByLocation() {
	city := c.Ctx.Input.Param(":city")
	countryCode := c.Ctx.Input.Param(":country")
	category := c.Ctx.Input.Param(":category")

	c.EventService.InvalidateCacheByLocation(
		city,
		countryCode,
		category,
	)

	c.Data["json"] = map[string]interface{}{
		"message": "event cache cleared",
	}

	c.ServeJSON()
}