package routers

import (
	"Event-Explorer/controllers"
	"Event-Explorer/services"

	beego "github.com/beego/beego/v2/server/web"
)

func Register(
	locationService *services.LocationService,
	eventService *services.EventService,
) {
	locationController := controllers.NewLocationController(
		locationService,
	)

	eventController := controllers.NewEventController(
		eventService,
	)

	beego.Router(
		"/",
		&controllers.MainController{},
	)

	beego.Router(
		"/api/locations/autocomplete",
		locationController,
		"get:Autocomplete",
	)

	beego.Router(
		"/api/locations/:placeId",
		locationController,
		"get:GetPlaceDetails",
	)

	beego.Router(
		"/events",
		eventController,
		"get:List",
	)

	beego.Router(
		"/events/:eventId",
		eventController,
		"get:Details",
	)

	beego.Router(
		"/redirect/:eventId",
		eventController,
		"get:Redirect",
	)
}
