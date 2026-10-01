package clients

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewGooglePlacesClient(t *testing.T) {
	client := NewGooglePlacesClient("test-api-key")

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

func TestGooglePlacesClient_Autocomplete(t *testing.T) {
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
				"suggestions": [
					{
						"placePrediction": {
							"placeId": "place-1",
							"text": {
								"text": "Toronto, ON, Canada"
							}
						}
					},
					{
						"placePrediction": {
							"placeId": "place-2",
							"text": {
								"text": "Toronto, OH, USA"
							}
						}
					}
				]
			}`,
			expectedCount: 2,
		},
		{
			name:          "empty suggestions",
			statusCode:    http.StatusOK,
			responseBody:  `{"suggestions":[]}`,
			expectedCount: 0,
		},
		{
			name:         "provider error",
			statusCode:   http.StatusBadRequest,
			responseBody: `{"error":"bad request"}`,
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
					if request.Method != http.MethodPost {
						t.Errorf(
							"expected POST, got %s",
							request.Method,
						)
					}

					if request.URL.Path != "/v1/places:autocomplete" {
						t.Errorf(
							"unexpected path: %s",
							request.URL.Path,
						)
					}

					if request.Header.Get("Content-Type") != "application/json" {
						t.Errorf("expected application/json content type")
					}

					if request.Header.Get("X-Goog-Api-Key") != "test-api-key" {
						t.Errorf("expected API key header")
					}

					var requestBody map[string]interface{}

					if err := json.NewDecoder(request.Body).Decode(
						&requestBody,
					); err != nil {
						t.Errorf(
							"failed to decode request body: %v",
							err,
						)
					}

					if requestBody["input"] != "Toronto" {
						t.Errorf(
							"expected input Toronto, got %v",
							requestBody["input"],
						)
					}

					if requestBody["sessionToken"] != "session-123" {
						t.Errorf(
							"expected session token session-123, got %v",
							requestBody["sessionToken"],
						)
					}

					if requestBody["languageCode"] != "en" {
						t.Errorf(
							"expected language en, got %v",
							requestBody["languageCode"],
						)
					}

					types, ok := requestBody["includedPrimaryTypes"].([]interface{})

					if !ok || len(types) != 1 || types[0] != "(cities)" {
						t.Errorf(
							"expected includedPrimaryTypes [(cities)], got %v",
							requestBody["includedPrimaryTypes"],
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

			client := NewGooglePlacesClient("test-api-key")

			oldBaseURL := googlePlacesBaseURL
			googlePlacesBaseURL = server.URL
			defer func() {
				googlePlacesBaseURL = oldBaseURL
			}()

			suggestions, err := client.Autocomplete(
				"Toronto",
				"session-123",
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

			if len(suggestions) != tt.expectedCount {
				t.Errorf(
					"expected %d suggestions, got %d",
					tt.expectedCount,
					len(suggestions),
				)
			}

			if tt.expectedCount > 0 {
				if suggestions[0].PlaceID != "place-1" {
					t.Errorf(
						"expected place ID place-1, got %s",
						suggestions[0].PlaceID,
					)
				}

				if suggestions[0].Text != "Toronto, ON, Canada" {
					t.Errorf(
						"unexpected suggestion text: %s",
						suggestions[0].Text,
					)
				}
			}
		})
	}
}

func TestGooglePlacesClient_AutocompleteRequestFailure(t *testing.T) {
	client := NewGooglePlacesClient("test-api-key")

	client.httpClient = &http.Client{
		Transport: roundTripperFunc(
			func(*http.Request) (*http.Response, error) {
				return nil, errors.New("network failure")
			},
		),
	}

	_, err := client.Autocomplete(
		"Toronto",
		"session-123",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGooglePlacesClient_GetPlaceDetails(t *testing.T) {
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
				"addressComponents": [
					{
						"longText": "Toronto",
						"shortText": "Toronto",
						"types": ["locality"]
					},
					{
						"longText": "Canada",
						"shortText": "CA",
						"types": ["country"]
					}
				]
			}`,
		},
		{
			name:         "place not found",
			statusCode:   http.StatusNotFound,
			responseBody: `{"error":"not found"}`,
			expectError:  true,
			checkError: func(t *testing.T, err error) {
				if !errors.Is(err, ErrGooglePlaceNotFound) {
					t.Errorf(
						"expected ErrGooglePlaceNotFound, got %v",
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
		{
			name:         "missing location data",
			statusCode:   http.StatusOK,
			responseBody: `{"addressComponents":[]}`,
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

					if request.URL.Path != "/v1/places/ChIJ123" {
						t.Errorf(
							"unexpected path: %s",
							request.URL.Path,
						)
					}

					if request.Header.Get("X-Goog-Api-Key") != "test-api-key" {
						t.Errorf("expected API key header")
					}

					if request.Header.Get("X-Goog-FieldMask") != "addressComponents" {
						t.Errorf(
							"expected addressComponents field mask",
						)
					}

					if request.URL.Query().Get("sessionToken") != "session-123" {
						t.Errorf(
							"expected session token",
						)
					}

					if request.URL.Query().Get("languageCode") != "en" {
						t.Errorf(
							"expected language en",
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

			client := NewGooglePlacesClient("test-api-key")

			oldBaseURL := googlePlacesBaseURL
			googlePlacesBaseURL = server.URL
			defer func() {
				googlePlacesBaseURL = oldBaseURL
			}()

			location, err := client.GetPlaceDetails(
				"ChIJ123",
				"session-123",
			)

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

			if location.City != "Toronto" {
				t.Errorf(
					"expected city Toronto, got %s",
					location.City,
				)
			}

			if location.CountryCode != "CA" {
				t.Errorf(
					"expected country code CA, got %s",
					location.CountryCode,
				)
			}
		})
	}
}

func TestGooglePlacesClient_GetPlaceDetailsRequestFailure(t *testing.T) {
	client := NewGooglePlacesClient("test-api-key")

	client.httpClient = &http.Client{
		Transport: roundTripperFunc(
			func(*http.Request) (*http.Response, error) {
				return nil, errors.New("network failure")
			},
		),
	}

	_, err := client.GetPlaceDetails(
		"ChIJ123",
		"session-123",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return f(request)
}
