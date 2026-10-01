package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"Event-Explorer/models"
	"Event-Explorer/services"

	"github.com/beego/beego/v2/server/web/context"
)

type mockLocationGoogleClient struct {
	autocompleteFunc func(
		string,
		string,
	) ([]models.LocationSuggestion, error)

	placeDetailsFunc func(
		string,
		string,
	) (*models.Location, error)
}

func (m *mockLocationGoogleClient) Autocomplete(
	input string,
	sessionToken string,
) ([]models.LocationSuggestion, error) {
	return m.autocompleteFunc(input, sessionToken)
}

func (m *mockLocationGoogleClient) GetPlaceDetails(
	placeID string,
	sessionToken string,
) (*models.Location, error) {
	return m.placeDetailsFunc(placeID, sessionToken)
}

func newTestLocationController(
	client *mockLocationGoogleClient,
) *LocationController {
	locationService := services.NewLocationService(client)

	return NewLocationController(locationService)
}

func prepareControllerContext(
	controller *LocationController,
	response *httptest.ResponseRecorder,
	request *http.Request,
) {
	ctx := context.NewContext()
	ctx.Reset(response, request)
	controller.Ctx = ctx
	controller.Data = make(map[interface{}]interface{})
}

func TestNewLocationController(t *testing.T) {
	client := &mockLocationGoogleClient{}

	controller := newTestLocationController(client)

	if controller == nil {
		t.Fatal("expected controller, got nil")
	}

	if controller.LocationService == nil {
		t.Fatal("expected location service, got nil")
	}
}

func TestLocationController_Autocomplete(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		sessionToken string
		setupMock    func() *mockLocationGoogleClient
		expectedCode int
	}{
		{
			name:         "success",
			input:        "Toronto",
			sessionToken: "session-123",
			setupMock: func() *mockLocationGoogleClient {
				return &mockLocationGoogleClient{
					autocompleteFunc: func(
						input string,
						sessionToken string,
					) ([]models.LocationSuggestion, error) {
						return []models.LocationSuggestion{
							{
								PlaceID: "place-1",
								Text:    "Toronto, ON, Canada",
							},
						}, nil
					},
				}
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid input",
			input:        "To",
			sessionToken: "session-123",
			setupMock: func() *mockLocationGoogleClient {
				return &mockLocationGoogleClient{
					autocompleteFunc: func(
						string,
						string,
					) ([]models.LocationSuggestion, error) {
						return nil, nil
					},
				}
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing session token",
			input:        "Toronto",
			sessionToken: "",
			setupMock: func() *mockLocationGoogleClient {
				return &mockLocationGoogleClient{
					autocompleteFunc: func(
						string,
						string,
					) ([]models.LocationSuggestion, error) {
						return nil, nil
					},
				}
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := newTestLocationController(
				tt.setupMock(),
			)

			request := httptest.NewRequest(
				http.MethodGet,
				"/api/locations/autocomplete?input="+
					tt.input+
					"&sessionToken="+
					tt.sessionToken,
				nil,
			)

			response := httptest.NewRecorder()

			prepareControllerContext(controller, response, request)

			runController(controller.Autocomplete)

			if response.Code != tt.expectedCode {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedCode,
					response.Code,
				)
			}
		})
	}
}

func TestLocationController_GetPlaceDetails(t *testing.T) {
	tests := []struct {
		name         string
		placeID      string
		sessionToken string
		setupMock    func() *mockLocationGoogleClient
		expectedCode int
	}{
		{
			name:         "success",
			placeID:      "place-123",
			sessionToken: "session-123",
			setupMock: func() *mockLocationGoogleClient {
				return &mockLocationGoogleClient{
					placeDetailsFunc: func(
						placeID string,
						sessionToken string,
					) (*models.Location, error) {
						return &models.Location{
							City:        "Toronto",
							CountryCode: "CA",
						}, nil
					},
				}
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "location not found",
			placeID:      "unknown-place",
			sessionToken: "session-123",
			setupMock: func() *mockLocationGoogleClient {
				return &mockLocationGoogleClient{
					placeDetailsFunc: func(
						placeID string,
						sessionToken string,
					) (*models.Location, error) {
						return nil, services.ErrLocationNotFound
					},
				}
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing place ID",
			placeID:      "",
			sessionToken: "session-123",
			setupMock: func() *mockLocationGoogleClient {
				return &mockLocationGoogleClient{
					placeDetailsFunc: func(
						placeID string,
						sessionToken string,
					) (*models.Location, error) {
						return nil, errors.New(
							"place ID is required",
						)
					},
				}
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing session token",
			placeID:      "place-123",
			sessionToken: "",
			setupMock: func() *mockLocationGoogleClient {
				return &mockLocationGoogleClient{
					placeDetailsFunc: func(
						placeID string,
						sessionToken string,
					) (*models.Location, error) {
						return nil, errors.New(
							"session token is required",
						)
					},
				}
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "provider failure",
			placeID:      "place-123",
			sessionToken: "session-123",
			setupMock: func() *mockLocationGoogleClient {
				return &mockLocationGoogleClient{
					placeDetailsFunc: func(
						placeID string,
						sessionToken string,
					) (*models.Location, error) {
						return nil, errors.New(
							"google places request failed",
						)
					},
				}
			},
			expectedCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := newTestLocationController(
				tt.setupMock(),
			)

			request := httptest.NewRequest(
				http.MethodGet,
				"/api/locations/"+tt.placeID+
					"?sessionToken="+tt.sessionToken,
				nil,
			)

			response := httptest.NewRecorder()

			prepareControllerContext(controller, response, request)

			controller.Ctx.Input.SetParam(
				":placeId",
				tt.placeID,
			)

			runController(controller.GetPlaceDetails)

			if response.Code != tt.expectedCode {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedCode,
					response.Code,
				)
			}
		})
	}
}
