package services

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"Event-Explorer/clients"
	"Event-Explorer/models"
)

type mockGooglePlacesClient struct {
	autocompleteFunc func(string, string) ([]models.LocationSuggestion, error)
	placeDetailsFunc func(string, string) (*models.Location, error)
}

func (m *mockGooglePlacesClient) Autocomplete(
	input string,
	sessionToken string,
) ([]models.LocationSuggestion, error) {
	return m.autocompleteFunc(input, sessionToken)
}

func (m *mockGooglePlacesClient) GetPlaceDetails(
	placeID string,
	sessionToken string,
) (*models.Location, error) {
	return m.placeDetailsFunc(placeID, sessionToken)
}

func TestLocationServiceAutocomplete(t *testing.T) {
	expectedSuggestions := []models.LocationSuggestion{
		{
			PlaceID: "mock-toronto",
			Text:    "Toronto, Canada",
		},
		{
			PlaceID: "mock-toronto-2",
			Text:    "Toronto, Canada",
		},
	}

	tests := []struct {
		name              string
		input             string
		sessionToken      string
		clientSuggestions []models.LocationSuggestion
		clientError       error
		want              []models.LocationSuggestion
		wantErr           bool
		wantErrMessage    string
		wantClientCalled  bool
		wantClientInput   string
		wantClientToken   string
	}{
		{
			name:              "successful autocomplete",
			input:             "Toronto",
			sessionToken:      "session-123",
			clientSuggestions: expectedSuggestions,
			want:              expectedSuggestions,
			wantClientCalled:  true,
			wantClientInput:   "Toronto",
			wantClientToken:   "session-123",
		},
		{
			name:              "trims input before calling client",
			input:             "  Toronto  ",
			sessionToken:      "session-123",
			clientSuggestions: expectedSuggestions,
			want:              expectedSuggestions,
			wantClientCalled:  true,
			wantClientInput:   "Toronto",
			wantClientToken:   "session-123",
		},
		{
			name:             "empty input",
			input:            "",
			sessionToken:     "session-123",
			wantErr:          true,
			wantErrMessage:   "input is required",
			wantClientCalled: false,
		},
		{
			name:             "whitespace input",
			input:            "   ",
			sessionToken:     "session-123",
			wantErr:          true,
			wantErrMessage:   "input is required",
			wantClientCalled: false,
		},
		{
			name:             "input shorter than three characters",
			input:            "To",
			sessionToken:     "session-123",
			wantErr:          true,
			wantErrMessage:   "input must contain at least 3 characters",
			wantClientCalled: false,
		},
		{
			name:              "input has exactly three characters",
			input:             "Tor",
			sessionToken:      "session-123",
			clientSuggestions: expectedSuggestions,
			want:              expectedSuggestions,
			wantClientCalled:  true,
			wantClientInput:   "Tor",
			wantClientToken:   "session-123",
		},
		{
			name:              "multibyte input with three characters",
			input:             "ঢাকা",
			sessionToken:      "session-123",
			clientSuggestions: expectedSuggestions,
			want:              expectedSuggestions,
			wantClientCalled:  true,
			wantClientInput:   "ঢাকা",
			wantClientToken:   "session-123",
		},
		{
			name:             "input exceeds 200 bytes",
			input:            strings.Repeat("a", 201),
			sessionToken:     "session-123",
			wantErr:          true,
			wantErrMessage:   "input must not exceed 200 bytes",
			wantClientCalled: false,
		},
		{
			name:              "input contains exactly 200 bytes",
			input:             strings.Repeat("a", 200),
			sessionToken:      "session-123",
			clientSuggestions: expectedSuggestions,
			want:              expectedSuggestions,
			wantClientCalled:  true,
			wantClientInput:   strings.Repeat("a", 200),
			wantClientToken:   "session-123",
		},
		{
			name:             "multibyte input exceeds 200 bytes",
			input:            strings.Repeat("ঢ", 67),
			sessionToken:     "session-123",
			wantErr:          true,
			wantErrMessage:   "input must not exceed 200 bytes",
			wantClientCalled: false,
		},
		{
			name:             "missing session token",
			input:            "Toronto",
			sessionToken:     "",
			wantErr:          true,
			wantErrMessage:   "session token is required",
			wantClientCalled: false,
		},
		{
			name:              "whitespace session token is accepted by current service",
			input:             "Toronto",
			sessionToken:      "   ",
			clientSuggestions: expectedSuggestions,
			want:              expectedSuggestions,
			wantClientCalled:  true,
			wantClientInput:   "Toronto",
			wantClientToken:   "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var receivedInput string
			var receivedToken string

			client := &mockGooglePlacesClient{
				autocompleteFunc: func(
					input string,
					sessionToken string,
				) ([]models.LocationSuggestion, error) {
					called = true
					receivedInput = input
					receivedToken = sessionToken

					return tt.clientSuggestions, tt.clientError
				},
				placeDetailsFunc: func(
					placeID string,
					sessionToken string,
				) (*models.Location, error) {
					return nil, nil
				},
			}

			service := NewLocationService(client)

			got, err := service.Autocomplete(
				tt.input,
				tt.sessionToken,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if err.Error() != tt.wantErrMessage {
					t.Fatalf(
						"expected error %q, got %q",
						tt.wantErrMessage,
						err.Error(),
					)
				}
			} else {
				if err != nil {
					t.Fatalf(
						"unexpected error: %v",
						err,
					)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf(
						"expected %+v, got %+v",
						tt.want,
						got,
					)
				}
			}

			if called != tt.wantClientCalled {
				t.Fatalf(
					"expected client called=%v, got %v",
					tt.wantClientCalled,
					called,
				)
			}

			if tt.wantClientCalled {
				if receivedInput != tt.wantClientInput {
					t.Fatalf(
						"expected client input %q, got %q",
						tt.wantClientInput,
						receivedInput,
					)
				}

				if receivedToken != tt.wantClientToken {
					t.Fatalf(
						"expected client token %q, got %q",
						tt.wantClientToken,
						receivedToken,
					)
				}
			}
		})
	}
}

func TestLocationServiceAutocompleteReturnsClientError(
	t *testing.T,
) {
	expectedErr := errors.New("google autocomplete failed")

	client := &mockGooglePlacesClient{
		autocompleteFunc: func(
			input string,
			sessionToken string,
		) ([]models.LocationSuggestion, error) {
			return nil, expectedErr
		},
		placeDetailsFunc: func(
			placeID string,
			sessionToken string,
		) (*models.Location, error) {
			return nil, nil
		},
	}

	service := NewLocationService(client)

	_, err := service.Autocomplete(
		"Toronto",
		"session-123",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected client error, got %v",
			err,
		)
	}
}

func TestLocationServiceGetPlaceDetails(t *testing.T) {
	expectedLocation := &models.Location{
		City:        "Toronto",
		CountryCode: "CA",
	}

	tests := []struct {
		name              string
		placeID           string
		sessionToken      string
		clientLocation    *models.Location
		clientError       error
		want              *models.Location
		wantErr           bool
		wantErrValue      error
		wantErrMessage    string
		wantClientCalled  bool
		wantClientPlaceID string
		wantClientToken   string
	}{
		{
			name:              "successful place details",
			placeID:           "mock-toronto",
			sessionToken:      "session-123",
			clientLocation:    expectedLocation,
			want:              expectedLocation,
			wantClientCalled:  true,
			wantClientPlaceID: "mock-toronto",
			wantClientToken:   "session-123",
		},
		{
			name:              "trims place ID",
			placeID:           "  mock-toronto  ",
			sessionToken:      "session-123",
			clientLocation:    expectedLocation,
			want:              expectedLocation,
			wantClientCalled:  true,
			wantClientPlaceID: "mock-toronto",
			wantClientToken:   "session-123",
		},
		{
			name:             "missing place ID",
			placeID:          "",
			sessionToken:     "session-123",
			wantErr:          true,
			wantErrMessage:   "place ID is required",
			wantClientCalled: false,
		},
		{
			name:             "whitespace place ID",
			placeID:          "   ",
			sessionToken:     "session-123",
			wantErr:          true,
			wantErrMessage:   "place ID is required",
			wantClientCalled: false,
		},
		{
			name:             "missing session token",
			placeID:          "mock-toronto",
			sessionToken:     "",
			wantErr:          true,
			wantErrMessage:   "session token is required",
			wantClientCalled: false,
		},
		{
			name:              "google place not found",
			placeID:           "unknown-place",
			sessionToken:      "session-123",
			clientError:       clients.ErrGooglePlaceNotFound,
			wantErr:           true,
			wantErrValue:      ErrLocationNotFound,
			wantClientCalled:  true,
			wantClientPlaceID: "unknown-place",
			wantClientToken:   "session-123",
		},
		{
			name:              "google provider error",
			placeID:           "mock-toronto",
			sessionToken:      "session-123",
			clientError:       errors.New("google api unavailable"),
			wantErr:           true,
			wantErrValue:      nil,
			wantErrMessage:    "google api unavailable",
			wantClientCalled:  true,
			wantClientPlaceID: "mock-toronto",
			wantClientToken:   "session-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var receivedPlaceID string
			var receivedToken string

			client := &mockGooglePlacesClient{
				autocompleteFunc: func(
					input string,
					sessionToken string,
				) ([]models.LocationSuggestion, error) {
					return nil, nil
				},
				placeDetailsFunc: func(
					placeID string,
					sessionToken string,
				) (*models.Location, error) {
					called = true
					receivedPlaceID = placeID
					receivedToken = sessionToken

					return tt.clientLocation, tt.clientError
				},
			}

			service := NewLocationService(client)

			got, err := service.GetPlaceDetails(
				tt.placeID,
				tt.sessionToken,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if tt.wantErrValue != nil {
					if !errors.Is(err, tt.wantErrValue) {
						t.Fatalf(
							"expected error %v, got %v",
							tt.wantErrValue,
							err,
						)
					}
				}

				if tt.wantErrMessage != "" &&
					err.Error() != tt.wantErrMessage {
					t.Fatalf(
						"expected error %q, got %q",
						tt.wantErrMessage,
						err.Error(),
					)
				}
			} else {
				if err != nil {
					t.Fatalf(
						"unexpected error: %v",
						err,
					)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf(
						"expected %+v, got %+v",
						tt.want,
						got,
					)
				}
			}

			if called != tt.wantClientCalled {
				t.Fatalf(
					"expected client called=%v, got %v",
					tt.wantClientCalled,
					called,
				)
			}

			if tt.wantClientCalled {
				if receivedPlaceID != tt.wantClientPlaceID {
					t.Fatalf(
						"expected place ID %q, got %q",
						tt.wantClientPlaceID,
						receivedPlaceID,
					)
				}

				if receivedToken != tt.wantClientToken {
					t.Fatalf(
						"expected token %q, got %q",
						tt.wantClientToken,
						receivedToken,
					)
				}
			}
		})
	}
}
