package services

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"Event-Explorer/clients"
	"Event-Explorer/models"
)


// Error variables for event service
var (
	ErrInvalidEventID          = errors.New("invalid event ID")
	ErrEventNotFound           = errors.New("event not found")
	ErrTicketURLMissing        = errors.New("ticket URL is missing")
	ErrUnsafeTicketURL         = errors.New("ticket URL is not allowed")
	ErrInvalidEventSearchInput = errors.New("invalid event search input")
)

type TicketmasterClient interface {
	GetEvents(
		city string,
		countryCode string,
		category string,
	) ([]models.Event, error)

	GetEvent(eventID string) (models.Event, error)
}

type EventService struct {
	ticketmasterClient TicketmasterClient

	cacheMutex sync.RWMutex
	cache      map[string][]models.Event
}

type eventResult struct {
	category  string
	events    []models.Event
	cacheHit  bool
	err       error
}

// CacheLocation represents a cached location with city, country code, and category.
type CacheLocation struct {
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
	Category    string `json:"category"`
}

// NewEventService creates a new instance of EventService with the provided TicketmasterClient.
func NewEventService(
	ticketmasterClient TicketmasterClient,
) *EventService {
	return &EventService{
		ticketmasterClient: ticketmasterClient,
		cache:              make(map[string][]models.Event),
	}
}

// GetEvents retrieves events based on the provided city and country code.
func (s *EventService) GetEvents(
	city string,
	countryCode string,
) (
	[]models.Event,
	[]models.Event,
	bool,
	bool,
	error,
) {
	city = strings.TrimSpace(city)
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))

	if city == "" {
		return nil, nil, false, false,
			fmt.Errorf(
				"%w: city is required",
				ErrInvalidEventSearchInput,
			)
	}

	if len(city) > 80 {
		return nil, nil, false, false,
			fmt.Errorf(
				"%w: city must not exceed 80 characters",
				ErrInvalidEventSearchInput,
			)
	}

	if !isValidCountryCode(countryCode) {
		return nil, nil, false, false,
			fmt.Errorf(
				"%w: country code must contain 2 letters",
				ErrInvalidEventSearchInput,
			)
	}

	results := make(chan eventResult, 2)

	go func() {
		events, cacheHit, err := s.getCategoryEvents(
			city,
			countryCode,
			"Music",
		)

		results <- eventResult{
			category: "Music",
			events:   events,
			cacheHit: cacheHit,
			err:      err,
		}
	}()

	go func() {
		events, cacheHit, err := s.getCategoryEvents(
			city,
			countryCode,
			"Sports",
		)

		results <- eventResult{
			category: "Sports",
			events:   events,
			cacheHit: cacheHit,
			err:      err,
		}
	}()

	var musicEvents []models.Event
	var sportsEvents []models.Event

	var musicCacheHit bool
	var sportsCacheHit bool

	var musicErr error
	var sportsErr error

	for i := 0; i < 2; i++ {
		result := <-results

		switch result.category {
		case "Music":
			musicEvents = result.events
			musicCacheHit = result.cacheHit
			musicErr = result.err

		case "Sports":
			sportsEvents = result.events
			sportsCacheHit = result.cacheHit
			sportsErr = result.err
		}
	}

	if musicErr != nil && sportsErr != nil {
		return nil, nil, false, false,
			fmt.Errorf(
				"failed to fetch music and sports events: music: %v; sports: %v",
				musicErr,
				sportsErr,
			)
	}

	return musicEvents, sportsEvents, musicCacheHit, sportsCacheHit, nil
}

// getCategoryEvents retrieves events for a specific category (Music or Sports) based on the provided city and country code.
func (s *EventService) getCategoryEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, bool, error) {
	key := cacheKey(city, countryCode, category)

	s.cacheMutex.RLock()

	cachedEvents, exists := s.cache[key]

	s.cacheMutex.RUnlock()

	if exists {
		fmt.Printf(
			"cache hit: city=%s country=%s category=%s\n",
			city,
			countryCode,
			category,
		)

		return cachedEvents, true, nil
	}

	events, err := s.ticketmasterClient.GetEvents(
		city,
		countryCode,
		category,
	)

	if err != nil {
		return nil, false, err
	}

	s.cacheMutex.Lock()

	s.cache[key] = events

	s.cacheMutex.Unlock()

	return events, false, nil
}

// GetEvent retrieves a specific event based on the provided event ID.
func (s *EventService) GetEvent(
	eventID string,
) (models.Event, error) {
	eventID = strings.TrimSpace(eventID)

	if eventID == "" {
		return models.Event{},
			ErrInvalidEventID
	}

	event, err := s.ticketmasterClient.GetEvent(eventID)

	if err != nil {
		if errors.Is(err, clients.ErrTicketmasterEventNotFound) {
			return models.Event{}, ErrEventNotFound
		}

		return models.Event{}, err
	}

	return event, nil
}

// GetTicketURL retrieves the ticket URL for a specific event based on the provided event ID.
func (s *EventService) GetTicketURL(
	eventID string,
) (string, error) {
	event, err := s.GetEvent(eventID)

	if err != nil {
		return "", err
	}

	if strings.TrimSpace(event.TicketURL) == "" {
		return "", ErrTicketURLMissing
	}

	parsedURL, err := url.Parse(event.TicketURL)
	if err != nil {
		return "", ErrUnsafeTicketURL
	}

	if parsedURL.Scheme != "https" {
		return "", ErrUnsafeTicketURL
	}

	if parsedURL.Hostname() != "www.ticketmaster.ca" {
		return "", ErrUnsafeTicketURL
	}

	return parsedURL.String(), nil
}

// GetCachedLocations retrieves cached locations based on a search query.
func (s *EventService) GetCachedLocations(
	search string,
) []CacheLocation {
	search = strings.ToLower(strings.TrimSpace(search))

	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	locations := make([]CacheLocation, 0)

	seen := make(map[string]bool)

	for key := range s.cache {
		parts := strings.Split(key, "|")

		if len(parts) != 3 {
			continue
		}

		city := parts[0]
		countryCode := parts[1]
		category := parts[2]

		if search != "" &&
			!strings.Contains(
				strings.ToLower(city),
				search,
			) {
			continue
		}

		locationKey := city + "|" + countryCode

		if seen[locationKey] {
			continue
		}

		seen[locationKey] = true

		locations = append(
			locations,
			CacheLocation{
				City:        city,
				CountryCode: countryCode,
				Category:    category,
			},
		)
	}

	return locations
}

// InvalidateCache clears all cached events in the EventService.
func (s *EventService) InvalidateCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache = make(map[string][]models.Event)
}

// InvalidateCacheByLocation clears cached events for a specific city, country code, and category in the EventService.
func (s *EventService) InvalidateCacheByLocation(
	city string,
	countryCode string,
	category string,
) {
	key := cacheKey(city, countryCode, category)

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	delete(s.cache, key)
}

// cacheKey generates a cache key based on the provided city, country code, and category.
func cacheKey(
	city string,
	countryCode string,
	category string,
) string {
	return strings.ToLower(
		strings.TrimSpace(city) +
			"|" +
			strings.TrimSpace(countryCode) +
			"|" +
			strings.TrimSpace(category),
	)
}

// isValidCountryCode checks if the provided value is a valid 2-letter country code.
func isValidCountryCode(value string) bool {
	if len(value) != 2 {
		return false
	}

	return value[0] >= 'A' && value[0] <= 'Z' &&
		value[1] >= 'A' && value[1] <= 'Z'
}

// NewEventServiceWithClient creates a new instance of EventService with the provided TicketmasterClient.
func NewEventServiceWithClient(
	ticketmasterClient TicketmasterClient,
) *EventService {
	return NewEventService(ticketmasterClient)
}