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

// GetEvents retrieves events for the specified city and country code, categorized into Music and Sports.
func (s *EventService) GetEvents(
	city string,
	countryCode string,
) ([]models.Event, []models.Event, error) {
	city = strings.TrimSpace(city)
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))

	if city == "" {
		return nil, nil, fmt.Errorf("city is required")
	}

	if len(countryCode) != 2 {
		return nil, nil, fmt.Errorf("country code must contain 2 letters")
	}

	type result struct {
		category string
		events   []models.Event
		err      error
	}

	results := make(chan result, 2)

	go func() {
		events, err := s.getCategoryEvents(
			city,
			countryCode,
			"Music",
		)

		results <- result{
			category: "Music",
			events:   events,
			err:      err,
		}
	}()

	go func() {
		events, err := s.getCategoryEvents(
			city,
			countryCode,
			"Sports",
		)

		results <- result{
			category: "Sports",
			events:   events,
			err:      err,
		}
	}()

	var music []models.Event
	var sports []models.Event

	var musicErr error
	var sportsErr error

	for i := 0; i < 2; i++ {
		result := <-results

		if result.category == "Music" {
			music = result.events
			musicErr = result.err
		} else {
			sports = result.events
			sportsErr = result.err
		}
	}

	if musicErr != nil && sportsErr != nil {
		return nil, nil, fmt.Errorf("failed to fetch events")
	}

	return music, sports, nil
}

// GetEvent retrieves a single event by its ID.
func (s *EventService) GetEvent(eventID string) (models.Event, error) {
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

// getCategoryEvents retrieves events for a specific category (Music or Sports) and caches the results.
func (s *EventService) getCategoryEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, error) {
	key := cacheKey(city, countryCode, category)

	s.cacheMutex.RLock()
	cachedEvents, found := s.cache[key]
	s.cacheMutex.RUnlock()

	if found {
		return cachedEvents, nil
	}

	events, err := s.ticketmasterClient.GetEvents(
		city,
		countryCode,
		category,
	)
	if err != nil {
		return nil, err
	}

	s.cacheMutex.Lock()
	s.cache[key] = events
	s.cacheMutex.Unlock()

	return events, nil
}

func cacheKey(
	city string,
	countryCode string,
	category string,
) string {
	return strings.ToLower(
		city + "|" + countryCode + "|" + category,
	)
}

func NewEventServiceWithClient(
	ticketmasterClient *clients.TicketmasterClient,
) *EventService {
	return NewEventService(ticketmasterClient)
}

func (s *EventService) GetTicketURL(eventID string) (string, error) {
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