package services

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"Event-Explorer/clients"
	"Event-Explorer/models"
)

type LocationService struct {
	googlePlacesClient *clients.GooglePlacesClient
}

func NewLocationService(
	googlePlacesClient *clients.GooglePlacesClient,
) *LocationService {
	return &LocationService{
		googlePlacesClient: googlePlacesClient,
	}
}

func (s *LocationService) Autocomplete(
	ctx context.Context,
	input string,
	sessionToken string,
) ([]models.LocationSuggestion, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, fmt.Errorf("input is required")
	}

	if utf8.RuneCountInString(input) < 3 {
		return nil, fmt.Errorf("input must contain at least 3 characters")
	}

	if len([]byte(input)) > 200 {
		return nil, fmt.Errorf("input must not exceed 200 bytes")
	}

	if sessionToken == "" {
		return nil, fmt.Errorf("session token is required")
	}

	return s.googlePlacesClient.Autocomplete(
		ctx,
		input,
		sessionToken,
	)
}

func (s *LocationService) GetPlaceDetails(
	ctx context.Context,
	placeID string,
	sessionToken string,
) (*models.Location, error) {
	placeID = strings.TrimSpace(placeID)

	if placeID == "" {
		return nil, fmt.Errorf("place ID is required")
	}

	if sessionToken == "" {
		return nil, fmt.Errorf("session token is required")
	}

	return s.googlePlacesClient.GetPlaceDetails(
		ctx,
		placeID,
		sessionToken,
	)
}
