package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"Event-Explorer/models"
)

const ticketmasterBaseURL = "https://app.ticketmaster.com/discovery/v2"

type TicketmasterClient struct {
	apiKey     string
	httpClient *http.Client
}

type ticketmasterEventsResponse struct {
	Embedded *struct {
		Events []ticketmasterEvent `json:"events"`
	} `json:"_embedded"`
}

type ticketmasterEvent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`

	Dates struct {
		Start struct {
			LocalDate string `json:"localDate"`
			LocalTime string `json:"localTime"`
		} `json:"start"`
		Timezone string `json:"timezone"`
	} `json:"dates"`

	Embedded *struct {
		Venues []struct {
			Name string `json:"name"`

			City struct {
				Name string `json:"name"`
			} `json:"city"`

			State struct {
				Name      string `json:"name"`
				StateCode string `json:"stateCode"`
			} `json:"state"`

			Country struct {
				Name        string `json:"name"`
				CountryCode string `json:"countryCode"`
			} `json:"country"`

			Address struct {
				Line1 string `json:"line1"`
			} `json:"address"`
		} `json:"venues"`
	} `json:"_embedded"`

	Info       string `json:"info"`
	PleaseNote string `json:"pleaseNote"`
	URL        string `json:"url"`
}

func NewTicketmasterClient(
	apiKey string,
) *TicketmasterClient {
	return &TicketmasterClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *TicketmasterClient) GetEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, error) {
	params := url.Values{}
	params.Set("apikey", c.apiKey)
	params.Set("city", city)
	params.Set("countryCode", countryCode)
	params.Set("classificationName", category)
	params.Set("size", "6")

	requestURL := ticketmasterBaseURL + "/events.json?" + params.Encode()

	response, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("ticketmaster request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"ticketmaster returned status %d",
			response.StatusCode,
		)
	}

	var data ticketmasterEventsResponse

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf(
			"failed to decode ticketmaster response: %w",
			err,
		)
	}

	if data.Embedded == nil {
		return []models.Event{}, nil
	}

	events := make([]models.Event, 0, len(data.Embedded.Events))

	for _, event := range data.Embedded.Events {
		events = append(
			events,
			mapTicketmasterEvent(event, category),
		)
	}

	return events, nil
}

func mapTicketmasterEvent(
	event ticketmasterEvent,
	category string,
) models.Event {
	result := models.Event{
		ID:          event.ID,
		Name:        event.Name,
		LocalDate:   event.Dates.Start.LocalDate,
		LocalTime:   event.Dates.Start.LocalTime,
		Timezone:    event.Dates.Timezone,
		Description: event.Info,
		TicketURL:   event.URL,
		Category:    category,
	}

	if result.Description == "" {
		result.Description = event.PleaseNote
	}

	if len(event.Images) > 0 {
		result.ImageURL = event.Images[0].URL
	}

	if event.Embedded != nil &&
		len(event.Embedded.Venues) > 0 {

		venue := event.Embedded.Venues[0]

		result.Venue = venue.Name
		result.City = venue.City.Name
		result.State = venue.State.Name
		result.Country = venue.Country.Name
		result.CountryCode = venue.Country.CountryCode
		result.Address = venue.Address.Line1
	}

	return result
}