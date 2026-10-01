package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"Event-Explorer/clients"
	"Event-Explorer/models"
	"Event-Explorer/services"

	"github.com/beego/beego/v2/server/web/context"
)

type mockTicketmasterClient struct {
	getEventsFunc func(string, string, string) ([]models.Event, error)
	getEventFunc  func(string) (models.Event, error)
}

func (m *mockTicketmasterClient) GetEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, error) {
	return m.getEventsFunc(city, countryCode, category)
}

func (m *mockTicketmasterClient) GetEvent(eventID string) (models.Event, error) {
	return m.getEventFunc(eventID)
}

func newTestEventController(client *mockTicketmasterClient) *EventController {
	return NewEventController(services.NewEventServiceWithClient(client))
}

func prepareEventContext(
	controller *EventController,
	response *httptest.ResponseRecorder,
	request *http.Request,
) {
	ctx := context.NewContext()
	ctx.Reset(response, request)
	controller.Ctx = ctx
	controller.Data = make(map[interface{}]interface{})
}

func runController(handler func()) {
	defer func() {
		_ = recover()
	}()

	handler()
}

func TestNewEventController(t *testing.T) {
	controller := newTestEventController(&mockTicketmasterClient{})

	if controller == nil || controller.EventService == nil {
		t.Fatal("expected event controller with event service")
	}
}

func TestEventController_List(t *testing.T) {
	tests := []struct {
		name         string
		city         string
		countryCode  string
		getEventsErr error
		expectedCode int
	}{
		{
			name:         "success",
			city:         "Toronto",
			countryCode:  "ca",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid search input",
			city:         "",
			countryCode:  "CA",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "provider failure",
			city:         "Toronto",
			countryCode:  "CA",
			getEventsErr: errors.New("ticketmaster unavailable"),
			expectedCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := newTestEventController(&mockTicketmasterClient{
				getEventsFunc: func(string, string, string) ([]models.Event, error) {
					if tt.getEventsErr != nil {
						return nil, tt.getEventsErr
					}
					return []models.Event{{ID: "event-1"}}, nil
				},
				getEventFunc: func(string) (models.Event, error) {
					return models.Event{}, nil
				},
			})
			request := httptest.NewRequest(
				http.MethodGet,
				"/events?city="+url.QueryEscape(tt.city)+"&countryCode="+url.QueryEscape(tt.countryCode),
				nil,
			)
			response := httptest.NewRecorder()
			prepareEventContext(controller, response, request)

			runController(controller.List)

			if response.Code != tt.expectedCode {
				t.Fatalf("expected status %d, got %d", tt.expectedCode, response.Code)
			}
			if tt.expectedCode == http.StatusOK {
				if controller.TplName != "listing.tpl" {
					t.Errorf("expected listing.tpl, got %q", controller.TplName)
				}
				if controller.Data["City"] != tt.city {
					t.Errorf("expected city %q, got %v", tt.city, controller.Data["City"])
				}
			}
		})
	}
}

func TestEventController_Details(t *testing.T) {
	tests := []struct {
		name         string
		eventID      string
		getEventErr  error
		expectedCode int
	}{
		{name: "success", eventID: "event-1", expectedCode: http.StatusOK},
		{name: "invalid ID", expectedCode: http.StatusBadRequest},
		{name: "not found", eventID: "missing", getEventErr: clients.ErrTicketmasterEventNotFound, expectedCode: http.StatusNotFound},
		{name: "provider failure", eventID: "event-1", getEventErr: errors.New("provider failure"), expectedCode: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := newTestEventController(&mockTicketmasterClient{
				getEventsFunc: func(string, string, string) ([]models.Event, error) { return nil, nil },
				getEventFunc: func(string) (models.Event, error) {
					if tt.getEventErr != nil {
						return models.Event{}, tt.getEventErr
					}
					return models.Event{ID: tt.eventID, Name: "Test event"}, nil
				},
			})
			request := httptest.NewRequest(http.MethodGet, "/events/"+tt.eventID, nil)
			response := httptest.NewRecorder()
			prepareEventContext(controller, response, request)
			controller.Ctx.Input.SetParam(":eventId", tt.eventID)

			runController(controller.Details)

			if response.Code != tt.expectedCode {
				t.Fatalf("expected status %d, got %d", tt.expectedCode, response.Code)
			}
			if tt.expectedCode == http.StatusOK && controller.TplName != "details.tpl" {
				t.Errorf("expected details.tpl, got %q", controller.TplName)
			}
		})
	}
}

func TestEventController_Redirect(t *testing.T) {
	tests := []struct {
		name         string
		eventID      string
		event        models.Event
		getEventErr  error
		expectedCode int
	}{
		{name: "success", eventID: "event-1", event: models.Event{TicketURL: "https://www.ticketmaster.ca/event-1"}, expectedCode: http.StatusFound},
		{name: "invalid ID", expectedCode: http.StatusBadRequest},
		{name: "not found", eventID: "missing", getEventErr: clients.ErrTicketmasterEventNotFound, expectedCode: http.StatusNotFound},
		{name: "missing URL", eventID: "event-1", expectedCode: http.StatusBadRequest},
		{name: "invalid URL", eventID: "event-1", event: models.Event{TicketURL: "://bad"}, expectedCode: http.StatusBadRequest},
		{name: "wrong host", eventID: "event-1", event: models.Event{TicketURL: "https://example.com/event-1"}, expectedCode: http.StatusBadRequest},
		{name: "provider failure", eventID: "event-1", getEventErr: errors.New("provider failure"), expectedCode: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := newTestEventController(&mockTicketmasterClient{
				getEventsFunc: func(string, string, string) ([]models.Event, error) { return nil, nil },
				getEventFunc:  func(string) (models.Event, error) { return tt.event, tt.getEventErr },
			})
			request := httptest.NewRequest(http.MethodGet, "/redirect/"+tt.eventID, nil)
			response := httptest.NewRecorder()
			prepareEventContext(controller, response, request)
			controller.Ctx.Input.SetParam(":eventId", tt.eventID)

			runController(controller.Redirect)

			if response.Code != tt.expectedCode {
				t.Fatalf("expected status %d, got %d", tt.expectedCode, response.Code)
			}
			if tt.expectedCode == http.StatusFound && response.Header().Get("Location") != tt.event.TicketURL {
				t.Errorf("expected redirect to %q, got %q", tt.event.TicketURL, response.Header().Get("Location"))
			}
		})
	}
}

func TestEventController_CacheEndpoints(t *testing.T) {
	controller := newTestEventController(&mockTicketmasterClient{
		getEventsFunc: func(string, string, string) ([]models.Event, error) {
			return []models.Event{{ID: "event-1"}}, nil
		},
		getEventFunc: func(string) (models.Event, error) { return models.Event{}, nil },
	})

	_, _, _, _, err := controller.EventService.GetEvents("Toronto", "CA")
	if err != nil {
		t.Fatalf("seed cache: %v", err)
	}

	t.Run("search short query", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/events/cache/locations?search=to", nil)
		response := httptest.NewRecorder()
		prepareEventContext(controller, response, request)
		runController(controller.SearchCachedCities)
		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}
	})

	t.Run("search cached cities", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/events/cache/locations?search=tor", nil)
		response := httptest.NewRecorder()
		prepareEventContext(controller, response, request)
		runController(controller.SearchCachedCities)
		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}
	})

	t.Run("invalidate all", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/events/cache", nil)
		response := httptest.NewRecorder()
		prepareEventContext(controller, response, request)
		runController(controller.InvalidateCache)
		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}
	})

	t.Run("invalidate by location", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/api/events/cache/Toronto/CA/Music", nil)
		response := httptest.NewRecorder()
		prepareEventContext(controller, response, request)
		controller.Ctx.Input.SetParam(":city", "Toronto")
		controller.Ctx.Input.SetParam(":country", "CA")
		controller.Ctx.Input.SetParam(":category", "Music")
		runController(controller.InvalidateCacheByLocation)
		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}
	})
}
