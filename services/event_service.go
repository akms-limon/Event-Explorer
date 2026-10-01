package services

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"Event-Explorer/clients"
	"Event-Explorer/models"
)

var (
	ErrInvalidEventID   = errors.New("invalid event ID")
	ErrEventNotFound    = errors.New("event not found")
	ErrTicketURLMissing = errors.New("ticket URL is missing")
	ErrUnsafeTicketURL  = errors.New("ticket URL is not allowed")
)

type TicketmasterClient interface {
	GetEvents(
		city string,
		countryCode string,
		category string,
	) ([]models.Event, error)

	GetEvent(
		eventID string,
	) (models.Event, error)
}

type EventService struct {
	ticketmasterClient TicketmasterClient

	cacheMutex sync.RWMutex
	cache      map[string][]models.Event
}

func NewEventService(
	ticketmasterClient TicketmasterClient,
) *EventService {
	return &EventService{
		ticketmasterClient: ticketmasterClient,
		cache:              make(map[string][]models.Event),
	}
}

type eventResult struct {
	category string
	events   []models.Event
	cacheHit bool
	err      error
}

// GetEvents retrieves Music and Sports events concurrently.
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
		return nil, nil, false, false, fmt.Errorf("city is required")
	}

	if len(city) > 80 {
		return nil, nil, false, false, fmt.Errorf(
			"city must not exceed 80 characters",
		)
	}

	if len(countryCode) != 2 {
		return nil, nil, false, false, fmt.Errorf(
			"country code must contain 2 letters",
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

	var music []models.Event
	var sports []models.Event

	var musicCacheHit bool
	var sportsCacheHit bool

	var musicErr error
	var sportsErr error

	for i := 0; i < 2; i++ {
		result := <-results

		if result.category == "Music" {
			music = result.events
			musicCacheHit = result.cacheHit
			musicErr = result.err
		} else {
			sports = result.events
			sportsCacheHit = result.cacheHit
			sportsErr = result.err
		}
	}

	if musicErr != nil && sportsErr != nil {
		return nil, nil, false, false, fmt.Errorf(
			"failed to fetch events",
		)
	}

	return music, sports, musicCacheHit, sportsCacheHit, nil
}

// GetEvent retrieves a single event by its ID.
func (s *EventService) GetEvent(
	eventID string,
) (models.Event, error) {
	eventID = strings.TrimSpace(eventID)

	if eventID == "" {
		return models.Event{}, ErrInvalidEventID
	}

	event, err := s.ticketmasterClient.GetEvent(eventID)
	if err != nil {
		if errors.Is(err, clients.ErrTicketmasterEventNotFound) {
			return models.Event{}, ErrEventNotFound
		}

		return models.Event{}, err
	}

	if event.ID == "" {
		return models.Event{}, ErrEventNotFound
	}

	return event, nil
}

func (s *EventService) getCategoryEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, bool, error) {
	key := cacheKey(city, countryCode, category)

	s.cacheMutex.RLock()
	cachedEvents, found := s.cache[key]
	s.cacheMutex.RUnlock()

	if found {
		log.Printf(
			"cache hit: city=%s country=%s category=%s",
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

func (s *EventService) InvalidateCache() {
	s.cacheMutex.Lock()
	s.cache = make(map[string][]models.Event)
	s.cacheMutex.Unlock()
}

func (s *EventService) InvalidateCacheByLocation(
	city string,
	countryCode string,
	category string,
) {
	key := cacheKey(
		city,
		countryCode,
		category,
	)

	s.cacheMutex.Lock()
	delete(s.cache, key)
	s.cacheMutex.Unlock()
}

func cacheKey(
	city string,
	countryCode string,
	category string,
) string {
	return strings.ToLower(
		strings.TrimSpace(city) +
			"|" +
			strings.ToUpper(strings.TrimSpace(countryCode)) +
			"|" +
			strings.TrimSpace(category),
	)
}

func NewEventServiceWithClient(
	ticketmasterClient *clients.TicketmasterClient,
) *EventService {
	return NewEventService(ticketmasterClient)
}

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

	ticketURL, err := url.Parse(event.TicketURL)
	if err != nil {
		return "", ErrUnsafeTicketURL
	}

	if ticketURL.Scheme != "https" {
		return "", ErrUnsafeTicketURL
	}

	if ticketURL.Host != "www.ticketmaster.ca" {
		return "", ErrUnsafeTicketURL
	}

	return ticketURL.String(), nil
}