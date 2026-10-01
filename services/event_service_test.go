package services

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"

	"Event-Explorer/clients"
	"Event-Explorer/models"
)

type mockTicketmasterClient struct {
	mu sync.Mutex

	getEventsFunc func(
		city string,
		countryCode string,
		category string,
	) ([]models.Event, error)

	getEventFunc func(
		eventID string,
	) (models.Event, error)

	getEventsCalls []string
	getEventCalls  []string
}

func (m *mockTicketmasterClient) GetEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, error) {
	m.mu.Lock()
	m.getEventsCalls = append(
		m.getEventsCalls,
		city+"|"+countryCode+"|"+category,
	)
	m.mu.Unlock()

	return m.getEventsFunc(
		city,
		countryCode,
		category,
	)
}

func (m *mockTicketmasterClient) GetEvent(
	eventID string,
) (models.Event, error) {
	m.mu.Lock()
	m.getEventCalls = append(
		m.getEventCalls,
		eventID,
	)
	m.mu.Unlock()

	return m.getEventFunc(eventID)
}

func TestEventServiceGetEventsValidation(t *testing.T) {
	tests := []struct {
		name        string
		city        string
		countryCode string
		wantMessage string
	}{
		{
			name:        "empty city",
			city:        "",
			countryCode: "CA",
			wantMessage: "city is required",
		},
		{
			name:        "whitespace city",
			city:        "   ",
			countryCode: "CA",
			wantMessage: "city is required",
		},
		{
			name:        "city exceeds 80 characters",
			city:        strings.Repeat("a", 81),
			countryCode: "CA",
			wantMessage: "city must not exceed 80 characters",
		},
		{
			name:        "country code too short",
			city:        "Toronto",
			countryCode: "C",
			wantMessage: "country code must contain 2 letters",
		},
		{
			name:        "country code too long",
			city:        "Toronto",
			countryCode: "CAN",
			wantMessage: "country code must contain 2 letters",
		},
		{
			name:        "country code contains number",
			city:        "Toronto",
			countryCode: "C1",
			wantMessage: "country code must contain 2 letters",
		},
		{
			name:        "country code contains special character",
			city:        "Toronto",
			countryCode: "C-",
			wantMessage: "country code must contain 2 letters",
		},
		{
			name:        "lowercase country code is accepted",
			city:        "Toronto",
			countryCode: "ca",
			wantMessage: "",
		},
		{
			name:        "city and country code are trimmed",
			city:        "  Toronto  ",
			countryCode: " ca ",
			wantMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newSuccessfulTicketmasterMock()

			service := NewEventService(client)

			_, _, _, _, err := service.GetEvents(
				tt.city,
				tt.countryCode,
			)

			if tt.wantMessage != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !errors.Is(
					err,
					ErrInvalidEventSearchInput,
				) {
					t.Fatalf(
						"expected ErrInvalidEventSearchInput, got %v",
						err,
					)
				}

				if !strings.Contains(
					err.Error(),
					tt.wantMessage,
				) {
					t.Fatalf(
						"expected error containing %q, got %q",
						tt.wantMessage,
						err.Error(),
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"unexpected error: %v",
					err,
				)
			}
		})
	}
}

func TestEventServiceGetEventsSuccess(t *testing.T) {
	musicEvents := []models.Event{
		{
			ID:   "music-1",
			Name: "Music Event",
		},
	}

	sportsEvents := []models.Event{
		{
			ID:   "sports-1",
			Name: "Sports Event",
		},
	}

	client := &mockTicketmasterClient{
		getEventsFunc: func(
			city string,
			countryCode string,
			category string,
		) ([]models.Event, error) {
			switch category {
			case "Music":
				return musicEvents, nil

			case "Sports":
				return sportsEvents, nil

			default:
				return nil, nil
			}
		},
		getEventFunc: func(
			eventID string,
		) (models.Event, error) {
			return models.Event{}, nil
		},
	}

	service := NewEventService(client)

	gotMusic, gotSports, musicHit, sportsHit, err :=
		service.GetEvents(
			"Toronto",
			"CA",
		)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if !reflect.DeepEqual(gotMusic, musicEvents) {
		t.Fatalf(
			"expected music events %+v, got %+v",
			musicEvents,
			gotMusic,
		)
	}

	if !reflect.DeepEqual(gotSports, sportsEvents) {
		t.Fatalf(
			"expected sports events %+v, got %+v",
			sportsEvents,
			gotSports,
		)
	}

	if musicHit {
		t.Fatal("expected music cache hit to be false")
	}

	if sportsHit {
		t.Fatal("expected sports cache hit to be false")
	}
}

func TestEventServiceGetEventsCacheHit(t *testing.T) {
	var callCount int

	client := &mockTicketmasterClient{
		getEventsFunc: func(
			city string,
			countryCode string,
			category string,
		) ([]models.Event, error) {
			callCount++

			return []models.Event{
				{
					ID:   category + "-1",
					Name: category + " Event",
				},
			}, nil
		},
		getEventFunc: func(
			eventID string,
		) (models.Event, error) {
			return models.Event{}, nil
		},
	}

	service := NewEventService(client)

	_, _, musicHit1, sportsHit1, err :=
		service.GetEvents(
			"Toronto",
			"CA",
		)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if musicHit1 || sportsHit1 {
		t.Fatal("first request should not be a cache hit")
	}

	_, _, musicHit2, sportsHit2, err :=
		service.GetEvents(
			"Toronto",
			"CA",
		)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if !musicHit2 {
		t.Fatal("expected music cache hit")
	}

	if !sportsHit2 {
		t.Fatal("expected sports cache hit")
	}

	if callCount != 2 {
		t.Fatalf(
			"expected 2 Ticketmaster calls, got %d",
			callCount,
		)
	}
}

func TestEventServiceGetEventsMusicFailure(t *testing.T) {
	musicErr := errors.New("music provider failed")

	client := &mockTicketmasterClient{
		getEventsFunc: func(
			city string,
			countryCode string,
			category string,
		) ([]models.Event, error) {
			if category == "Music" {
				return nil, musicErr
			}

			return []models.Event{
				{
					ID:   "sports-1",
					Name: "Sports Event",
				},
			}, nil
		},
		getEventFunc: func(
			eventID string,
		) (models.Event, error) {
			return models.Event{}, nil
		},
	}

	service := NewEventService(client)

	music, sports, musicHit, sportsHit, err :=
		service.GetEvents(
			"Toronto",
			"CA",
		)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if music != nil {
		t.Fatalf("expected no music events, got %+v", music)
	}

	if len(sports) != 1 {
		t.Fatalf(
			"expected one sports event, got %d",
			len(sports),
		)
	}

	if musicHit || sportsHit {
		t.Fatal("expected both cache hit flags to be false")
	}
}

func TestEventServiceGetEventsSportsFailure(t *testing.T) {
	sportsErr := errors.New("sports provider failed")

	client := &mockTicketmasterClient{
		getEventsFunc: func(
			city string,
			countryCode string,
			category string,
		) ([]models.Event, error) {
			if category == "Sports" {
				return nil, sportsErr
			}

			return []models.Event{
				{
					ID:   "music-1",
					Name: "Music Event",
				},
			}, nil
		},
		getEventFunc: func(
			eventID string,
		) (models.Event, error) {
			return models.Event{}, nil
		},
	}

	service := NewEventService(client)

	music, sports, musicHit, sportsHit, err :=
		service.GetEvents(
			"Toronto",
			"CA",
		)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if len(music) != 1 {
		t.Fatalf(
			"expected one music event, got %d",
			len(music),
		)
	}

	if sports != nil {
		t.Fatalf(
			"expected no sports events, got %+v",
			sports,
		)
	}

	if musicHit || sportsHit {
		t.Fatal("expected both cache hit flags to be false")
	}
}

func TestEventServiceGetEventsBothCategoriesFail(t *testing.T) {
	client := &mockTicketmasterClient{
		getEventsFunc: func(
			city string,
			countryCode string,
			category string,
		) ([]models.Event, error) {
			return nil, errors.New(
				category + " provider failed",
			)
		},
		getEventFunc: func(
			eventID string,
		) (models.Event, error) {
			return models.Event{}, nil
		},
	}

	service := NewEventService(client)

	music, sports, musicHit, sportsHit, err :=
		service.GetEvents(
			"Toronto",
			"CA",
		)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"failed to fetch music and sports events",
	) {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if music != nil || sports != nil {
		t.Fatal("expected nil event lists")
	}

	if musicHit || sportsHit {
		t.Fatal("expected both cache hit flags to be false")
	}
}

func TestEventServiceGetEvent(t *testing.T) {
	expectedEvent := models.Event{
		ID:   "event-1",
		Name: "Test Event",
	}

	tests := []struct {
		name        string
		eventID     string
		clientEvent models.Event
		clientError error
		want        models.Event
		wantErr     error
	}{
		{
			name:        "successful event",
			eventID:     "event-1",
			clientEvent: expectedEvent,
			want:        expectedEvent,
		},
		{
			name:    "empty event ID",
			eventID: "",
			wantErr: ErrInvalidEventID,
		},
		{
			name:    "whitespace event ID",
			eventID: "   ",
			wantErr: ErrInvalidEventID,
		},
		{
			name:        "event not found",
			eventID:     "missing",
			clientError: clients.ErrTicketmasterEventNotFound,
			wantErr:     ErrEventNotFound,
		},
		{
			name:        "ticketmaster provider error",
			eventID:     "event-1",
			clientError: errors.New("ticketmaster unavailable"),
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockTicketmasterClient{
				getEventsFunc: func(
					city string,
					countryCode string,
					category string,
				) ([]models.Event, error) {
					return nil, nil
				},
				getEventFunc: func(
					eventID string,
				) (models.Event, error) {
					return tt.clientEvent, tt.clientError
				},
			}

			service := NewEventService(client)

			got, err := service.GetEvent(tt.eventID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf(
						"expected error %v, got %v",
						tt.wantErr,
						err,
					)
				}
				return
			}

			if tt.clientError != nil {
				if err == nil {
					t.Fatal("expected provider error")
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf(
					"expected %+v, got %+v",
					tt.want,
					got,
				)
			}
		})
	}
}

func TestEventServiceGetTicketURL(t *testing.T) {
	tests := []struct {
		name        string
		event       models.Event
		clientError error
		want        string
		wantErr     error
	}{
		{
			name: "valid ticketmaster URL",
			event: models.Event{
				ID:        "event-1",
				TicketURL: "https://www.ticketmaster.ca/event-1",
			},
			want: "https://www.ticketmaster.ca/event-1",
		},
		{
			name: "missing ticket URL",
			event: models.Event{
				ID:        "event-1",
				TicketURL: "",
			},
			wantErr: ErrTicketURLMissing,
		},
		{
			name: "whitespace ticket URL",
			event: models.Event{
				ID:        "event-1",
				TicketURL: "   ",
			},
			wantErr: ErrTicketURLMissing,
		},
		{
			name: "http URL",
			event: models.Event{
				ID:        "event-1",
				TicketURL: "http://www.ticketmaster.ca/event-1",
			},
			wantErr: ErrUnsafeTicketURL,
		},
		{
			name: "wrong hostname",
			event: models.Event{
				ID:        "event-1",
				TicketURL: "https://example.com/event-1",
			},
			wantErr: ErrUnsafeTicketURL,
		},
		{
			name: "ticketmaster.com hostname",
			event: models.Event{
				ID:        "event-1",
				TicketURL: "https://www.ticketmaster.com/event-1",
			},
			wantErr: ErrUnsafeTicketURL,
		},
		{
			name: "invalid URL",
			event: models.Event{
				ID:        "event-1",
				TicketURL: "://invalid-url",
			},
			wantErr: ErrUnsafeTicketURL,
		},
		{
			name:        "event not found",
			event:       models.Event{},
			clientError: clients.ErrTicketmasterEventNotFound,
			wantErr:     ErrEventNotFound,
		},
		{
			name:    "invalid event ID",
			event:   models.Event{},
			wantErr: ErrInvalidEventID,
		},
		{
			name:        "provider error",
			event:       models.Event{},
			clientError: errors.New("provider unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventID := tt.event.ID

			if eventID == "" &&
				tt.name != "invalid event ID" {
				eventID = "event-1"
			}

			client := &mockTicketmasterClient{
				getEventsFunc: func(
					city string,
					countryCode string,
					category string,
				) ([]models.Event, error) {
					return nil, nil
				},
				getEventFunc: func(
					id string,
				) (models.Event, error) {
					return tt.event, tt.clientError
				},
			}

			service := NewEventService(client)

			got, err := service.GetTicketURL(eventID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf(
						"expected error %v, got %v",
						tt.wantErr,
						err,
					)
				}
				return
			}

			if tt.clientError != nil {
				if err == nil {
					t.Fatal("expected provider error")
				}
				return
			}

			if got != tt.want {
				t.Fatalf(
					"expected URL %q, got %q",
					tt.want,
					got,
				)
			}
		})
	}
}

func TestEventServiceGetCachedLocations(t *testing.T) {
	client := newSuccessfulTicketmasterMock()

	service := NewEventService(client)

	service.cache[cacheKey(
		"Toronto",
		"CA",
		"Music",
	)] = []models.Event{}

	service.cache[cacheKey(
		"Toronto",
		"CA",
		"Sports",
	)] = []models.Event{}

	service.cache[cacheKey(
		"London",
		"GB",
		"Music",
	)] = []models.Event{}

	tests := []struct {
		name   string
		search string
		want   int
	}{
		{
			name:   "all cached locations",
			search: "",
			want:   2,
		},
		{
			name:   "search Toronto",
			search: "Toronto",
			want:   1,
		},
		{
			name:   "search case insensitive",
			search: "TOR",
			want:   1,
		},
		{
			name:   "search with whitespace",
			search: "  london  ",
			want:   1,
		},
		{
			name:   "no matching location",
			search: "Dhaka",
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.GetCachedLocations(tt.search)

			if len(got) != tt.want {
				t.Fatalf(
					"expected %d locations, got %d",
					tt.want,
					len(got),
				)
			}
		})
	}
}

func TestEventServiceInvalidateCache(t *testing.T) {
	client := newSuccessfulTicketmasterMock()

	service := NewEventService(client)

	service.cache["toronto|ca|music"] = []models.Event{
		{
			ID: "event-1",
		},
	}

	service.cache["london|gb|sports"] = []models.Event{
		{
			ID: "event-2",
		},
	}

	service.InvalidateCache()

	if len(service.cache) != 0 {
		t.Fatalf(
			"expected empty cache, got %d entries",
			len(service.cache),
		)
	}
}

func TestEventServiceInvalidateCacheByLocation(t *testing.T) {
	client := newSuccessfulTicketmasterMock()

	service := NewEventService(client)

	service.cache[cacheKey(
		"Toronto",
		"CA",
		"Music",
	)] = []models.Event{
		{
			ID: "music-1",
		},
	}

	service.cache[cacheKey(
		"Toronto",
		"CA",
		"Sports",
	)] = []models.Event{
		{
			ID: "sports-1",
		},
	}

	service.InvalidateCacheByLocation(
		" Toronto ",
		"CA",
		"Music",
	)

	if _, exists := service.cache[cacheKey("Toronto", "CA", "Music")]; exists {
		t.Fatal("expected music cache to be removed")
	}

	if _, exists := service.cache[cacheKey("Toronto", "CA", "Sports")]; !exists {
		t.Fatal("expected sports cache to remain")
	}
}

func TestCacheKey(t *testing.T) {
	tests := []struct {
		name        string
		city        string
		countryCode string
		category    string
		want        string
	}{
		{
			name:        "normal values",
			city:        "Toronto",
			countryCode: "CA",
			category:    "Music",
			want:        "toronto|ca|music",
		},
		{
			name:        "trims and lowercases values",
			city:        "  Toronto  ",
			countryCode: " CA ",
			category:    " MUSIC ",
			want:        "toronto|ca|music",
		},
		{
			name:        "empty values",
			city:        "",
			countryCode: "",
			category:    "",
			want:        "||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cacheKey(
				tt.city,
				tt.countryCode,
				tt.category,
			)

			if got != tt.want {
				t.Fatalf(
					"expected %q, got %q",
					tt.want,
					got,
				)
			}
		})
	}
}

func TestIsValidCountryCode(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "valid uppercase code",
			value: "CA",
			want:  true,
		},
		{
			name:  "valid another code",
			value: "GB",
			want:  true,
		},
		{
			name:  "lowercase code",
			value: "ca",
			want:  false,
		},
		{
			name:  "one character",
			value: "C",
			want:  false,
		},
		{
			name:  "three characters",
			value: "CAN",
			want:  false,
		},
		{
			name:  "number",
			value: "C1",
			want:  false,
		},
		{
			name:  "special character",
			value: "C-",
			want:  false,
		},
		{
			name:  "empty",
			value: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidCountryCode(tt.value)

			if got != tt.want {
				t.Fatalf(
					"expected %v, got %v",
					tt.want,
					got,
				)
			}
		})
	}
}

func TestNewEventServiceWithClient(t *testing.T) {
	client := newSuccessfulTicketmasterMock()

	service := NewEventServiceWithClient(client)

	if service == nil {
		t.Fatal("expected service, got nil")
	}

	if service.ticketmasterClient != client {
		t.Fatal("expected provided client to be assigned")
	}

	if service.cache == nil {
		t.Fatal("expected cache to be initialized")
	}
}

func newSuccessfulTicketmasterMock() *mockTicketmasterClient {
	return &mockTicketmasterClient{
		getEventsFunc: func(
			city string,
			countryCode string,
			category string,
		) ([]models.Event, error) {
			return []models.Event{
				{
					ID:   category + "-1",
					Name: category + " Event",
				},
			}, nil
		},
		getEventFunc: func(
			eventID string,
		) (models.Event, error) {
			return models.Event{
				ID:        eventID,
				Name:      "Test Event",
				TicketURL: "https://www.ticketmaster.ca/" + eventID,
			}, nil
		},
	}
}
