package clients

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTicketmasterClient(t *testing.T) {
	client := NewTicketmasterClient("test-api-key")

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	if client.apiKey != "test-api-key" {
		t.Errorf(
			"expected api key %q, got %q",
			"test-api-key",
			client.apiKey,
		)
	}

	if client.httpClient == nil {
		t.Fatal("expected http client, got nil")
	}

	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf(
			"expected timeout %v, got %v",
			10*time.Second,
			client.httpClient.Timeout,
		)
	}
}

func TestTicketmasterClient_GetEvents(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectError   bool
		expectedCount int
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			responseBody: `{
				"_embedded": {
					"events": [
						{
							"id": "event-1",
							"name": "Toronto Music Festival",
							"images": [
								{
									"url": "https://example.com/image.jpg"
								}
							],
							"classifications": [
								{
									"primary": true,
									"segment": {
										"name": "Music"
									}
								}
							],
							"dates": {
								"start": {
									"localDate": "2026-10-10",
									"localTime": "19:00:00"
								},
								"timezone": "America/Toronto"
							},
							"_embedded": {
								"venues": [
									{
										"name": "Scotiabank Arena",
										"city": {
											"name": "Toronto"
										},
										"state": {
											"name": "Ontario",
											"stateCode": "ON"
										},
										"country": {
											"name": "Canada",
											"countryCode": "CA"
										},
										"address": {
											"line1": "40 Bay Street"
										}
									}
								]
							},
							"info": "Music festival information",
							"pleaseNote": "Please arrive early",
							"url": "https://www.ticketmaster.ca/event-1"
						}
					]
				}
			}`,
			expectedCount: 1,
		},
		{
			name:          "empty response",
			statusCode:    http.StatusOK,
			responseBody:  `{}`,
			expectedCount: 0,
		},
		{
			name:         "provider error",
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"error":"server error"}`,
			expectError:  true,
		},
		{
			name:         "invalid json",
			statusCode:   http.StatusOK,
			responseBody: `invalid-json`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(
					writer http.ResponseWriter,
					request *http.Request,
				) {
					if request.Method != http.MethodGet {
						t.Errorf(
							"expected GET, got %s",
							request.Method,
						)
					}

					if request.URL.Path != "/discovery/v2/events.json" {
						t.Errorf(
							"unexpected path: %s",
							request.URL.Path,
						)
					}

					query := request.URL.Query()

					if query.Get("apikey") != "test-api-key" {
						t.Errorf(
							"expected api key, got %q",
							query.Get("apikey"),
						)
					}

					if query.Get("city") != "Toronto" {
						t.Errorf(
							"expected city Toronto, got %q",
							query.Get("city"),
						)
					}

					if query.Get("countryCode") != "CA" {
						t.Errorf(
							"expected country code CA, got %q",
							query.Get("countryCode"),
						)
					}

					if query.Get("classificationName") != "Music" {
						t.Errorf(
							"expected category Music, got %q",
							query.Get("classificationName"),
						)
					}

					if query.Get("size") != "6" {
						t.Errorf(
							"expected size 6, got %q",
							query.Get("size"),
						)
					}

					writer.Header().Set(
						"Content-Type",
						"application/json",
					)

					writer.WriteHeader(tt.statusCode)

					_, _ = writer.Write(
						[]byte(tt.responseBody),
					)
				}),
			)

			defer server.Close()

			client := NewTicketmasterClient("test-api-key")

			oldBaseURL := ticketmasterBaseURL
			ticketmasterBaseURL = server.URL + "/discovery/v2"

			defer func() {
				ticketmasterBaseURL = oldBaseURL
			}()

			events, err := client.GetEvents(
				"Toronto",
				"CA",
				"Music",
			)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"expected no error, got %v",
					err,
				)
			}

			if len(events) != tt.expectedCount {
				t.Errorf(
					"expected %d events, got %d",
					tt.expectedCount,
					len(events),
				)
			}

			if tt.expectedCount == 1 {
				event := events[0]

				if event.ID != "event-1" {
					t.Errorf(
						"expected event ID event-1, got %s",
						event.ID,
					)
				}

				if event.Name != "Toronto Music Festival" {
					t.Errorf(
						"unexpected event name: %s",
						event.Name,
					)
				}

				if event.Category != "Music" {
					t.Errorf(
						"expected category Music, got %s",
						event.Category,
					)
				}

				if event.ImageURL != "https://example.com/image.jpg" {
					t.Errorf(
						"unexpected image URL: %s",
						event.ImageURL,
					)
				}

				if event.Venue != "Scotiabank Arena" {
					t.Errorf(
						"unexpected venue: %s",
						event.Venue,
					)
				}

				if event.City != "Toronto" {
					t.Errorf(
						"unexpected city: %s",
						event.City,
					)
				}

				if event.State != "Ontario" {
					t.Errorf(
						"unexpected state: %s",
						event.State,
					)
				}

				if event.Country != "Canada" {
					t.Errorf(
						"unexpected country: %s",
						event.Country,
					)
				}

				if event.CountryCode != "CA" {
					t.Errorf(
						"unexpected country code: %s",
						event.CountryCode,
					)
				}

				if event.Address != "40 Bay Street" {
					t.Errorf(
						"unexpected address: %s",
						event.Address,
					)
				}
			}
		})
	}
}

func TestTicketmasterClient_GetEventsRequestFailure(t *testing.T) {
	client := NewTicketmasterClient("test-api-key")

	client.httpClient = &http.Client{
		Transport: roundTripperFunc(
			func(*http.Request) (*http.Response, error) {
				return nil, errors.New("network failure")
			},
		),
	}

	_, err := client.GetEvents(
		"Toronto",
		"CA",
		"Music",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTicketmasterClient_GetEvent(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		expectError  bool
		checkError   func(t *testing.T, err error)
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			responseBody: `{
				"id": "event-123",
				"name": "Toronto Concert",
				"classifications": [
					{
						"primary": true,
						"segment": {
							"name": "Music"
						}
					}
				],
				"dates": {
					"start": {
						"localDate": "2026-11-10",
						"localTime": "20:00:00"
					},
					"timezone": "America/Toronto"
				},
				"info": "Concert information",
				"url": "https://www.ticketmaster.ca/event-123"
			}`,
		},
		{
			name:         "event not found",
			statusCode:   http.StatusNotFound,
			responseBody: `{"error":"not found"}`,
			expectError:  true,
			checkError: func(t *testing.T, err error) {
				if !errors.Is(err, ErrTicketmasterEventNotFound) {
					t.Errorf(
						"expected ErrTicketmasterEventNotFound, got %v",
						err,
					)
				}
			},
		},
		{
			name:         "provider error",
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"error":"server error"}`,
			expectError:  true,
		},
		{
			name:         "invalid json",
			statusCode:   http.StatusOK,
			responseBody: `invalid-json`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(
					writer http.ResponseWriter,
					request *http.Request,
				) {
					if request.Method != http.MethodGet {
						t.Errorf(
							"expected GET, got %s",
							request.Method,
						)
					}

					if request.URL.Path != "/discovery/v2/events/event-123.json" {
						t.Errorf(
							"unexpected path: %s",
							request.URL.Path,
						)
					}

					if request.URL.Query().Get("apikey") != "test-api-key" {
						t.Errorf("expected API key")
					}

					writer.Header().Set(
						"Content-Type",
						"application/json",
					)

					writer.WriteHeader(tt.statusCode)

					_, _ = writer.Write(
						[]byte(tt.responseBody),
					)
				}),
			)

			defer server.Close()

			client := NewTicketmasterClient("test-api-key")

			oldBaseURL := ticketmasterBaseURL
			ticketmasterBaseURL = server.URL + "/discovery/v2"

			defer func() {
				ticketmasterBaseURL = oldBaseURL
			}()

			event, err := client.GetEvent("event-123")

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if tt.checkError != nil {
					tt.checkError(t, err)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"expected no error, got %v",
					err,
				)
			}

			if event.ID != "event-123" {
				t.Errorf(
					"expected event ID event-123, got %s",
					event.ID,
				)
			}

			if event.Name != "Toronto Concert" {
				t.Errorf(
					"expected Toronto Concert, got %s",
					event.Name,
				)
			}

			if event.Category != "Music" {
				t.Errorf(
					"expected Music, got %s",
					event.Category,
				)
			}
		})
	}
}

func TestTicketmasterClient_GetEventRequestFailure(t *testing.T) {
	client := NewTicketmasterClient("test-api-key")

	client.httpClient = &http.Client{
		Transport: roundTripperFunc(
			func(*http.Request) (*http.Response, error) {
				return nil, errors.New("network failure")
			},
		),
	}

	_, err := client.GetEvent("event-123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMapTicketmasterEvent(t *testing.T) {
	t.Run("explicit category and info description", func(t *testing.T) {
		event := ticketmasterEvent{
			ID:   "event-1",
			Name: "Test Event",
			Images: []struct {
				URL string `json:"url"`
			}{
				{
					URL: "https://example.com/image.jpg",
				},
			},
			Dates: struct {
				Start struct {
					LocalDate string `json:"localDate"`
					LocalTime string `json:"localTime"`
				} `json:"start"`
				Timezone string `json:"timezone"`
			}{
				Start: struct {
					LocalDate string `json:"localDate"`
					LocalTime string `json:"localTime"`
				}{
					LocalDate: "2026-10-10",
					LocalTime: "19:00:00",
				},
				Timezone: "America/Toronto",
			},
			Info:       "Event info",
			PleaseNote: "Please note",
			URL:        "https://www.ticketmaster.ca/event-1",
		}

		result := mapTicketmasterEvent(event, "Music")

		if result.Category != "Music" {
			t.Errorf(
				"expected Music, got %s",
				result.Category,
			)
		}

		if result.Description != "Event info" {
			t.Errorf(
				"expected Event info, got %s",
				result.Description,
			)
		}

		if result.ImageURL != "https://example.com/image.jpg" {
			t.Errorf(
				"unexpected image URL: %s",
				result.ImageURL,
			)
		}
	})

	t.Run("primary classification", func(t *testing.T) {
		event := ticketmasterEvent{
			Classifications: []struct {
				Primary bool `json:"primary"`
				Segment struct {
					Name string `json:"name"`
				} `json:"segment"`
			}{
				{
					Primary: true,
					Segment: struct {
						Name string `json:"name"`
					}{
						Name: "Sports",
					},
				},
			},
		}

		result := mapTicketmasterEvent(event, "")

		if result.Category != "Sports" {
			t.Errorf(
				"expected Sports, got %s",
				result.Category,
			)
		}
	})

	t.Run("first classification fallback", func(t *testing.T) {
		event := ticketmasterEvent{
			Classifications: []struct {
				Primary bool `json:"primary"`
				Segment struct {
					Name string `json:"name"`
				} `json:"segment"`
			}{
				{
					Primary: false,
					Segment: struct {
						Name string `json:"name"`
					}{
						Name: "Music",
					},
				},
			},
		}

		result := mapTicketmasterEvent(event, "")

		if result.Category != "Music" {
			t.Errorf(
				"expected Music, got %s",
				result.Category,
			)
		}
	})

	t.Run("please note fallback", func(t *testing.T) {
		event := ticketmasterEvent{
			PleaseNote: "Please arrive early",
		}

		result := mapTicketmasterEvent(event, "")

		if result.Description != "Please arrive early" {
			t.Errorf(
				"expected Please arrive early, got %s",
				result.Description,
			)
		}
	})

	t.Run("venue mapping", func(t *testing.T) {
		event := ticketmasterEvent{
			Embedded: &struct {
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
			}{
				Venues: []struct {
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
				}{
					{
						Name: "Scotiabank Arena",
						City: struct {
							Name string `json:"name"`
						}{
							Name: "Toronto",
						},
						State: struct {
							Name      string `json:"name"`
							StateCode string `json:"stateCode"`
						}{
							Name: "Ontario",
						},
						Country: struct {
							Name        string `json:"name"`
							CountryCode string `json:"countryCode"`
						}{
							Name:        "Canada",
							CountryCode: "CA",
						},
						Address: struct {
							Line1 string `json:"line1"`
						}{
							Line1: "40 Bay Street",
						},
					},
				},
			},
		}

		result := mapTicketmasterEvent(event, "")

		if result.Venue != "Scotiabank Arena" {
			t.Errorf(
				"expected Scotiabank Arena, got %s",
				result.Venue,
			)
		}

		if result.City != "Toronto" {
			t.Errorf(
				"expected Toronto, got %s",
				result.City,
			)
		}

		if result.CountryCode != "CA" {
			t.Errorf(
				"expected CA, got %s",
				result.CountryCode,
			)
		}

		if result.Address != "40 Bay Street" {
			t.Errorf(
				"expected 40 Bay Street, got %s",
				result.Address,
			)
		}
	})
}
