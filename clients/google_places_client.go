package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"Event-Explorer/models"
)

const googlePlacesBaseURL = "https://places.googleapis.com"

type GooglePlacesClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewGooglePlacesClient(apiKey string) *GooglePlacesClient {
	return &GooglePlacesClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type autocompleteRequest struct {
	Input                string   `json:"input"`
	IncludedPrimaryTypes []string `json:"includedPrimaryTypes"`
	SessionToken         string   `json:"sessionToken"`
	LanguageCode         string   `json:"languageCode"`
}

type autocompleteResponse struct {
	Suggestions []struct {
		PlacePrediction struct {
			PlaceID string `json:"placeId"`
			Text    struct {
				Text string `json:"text"`
			} `json:"text"`
		} `json:"placePrediction"`
	} `json:"suggestions"`
}

type placeDetailsResponse struct {
	AddressComponents []struct {
		LongText  string   `json:"longText"`
		ShortText string   `json:"shortText"`
		Types     []string `json:"types"`
	} `json:"addressComponents"`
}

func (c *GooglePlacesClient) Autocomplete(
	input string,
	sessionToken string,
) ([]models.LocationSuggestion, error) {
	requestBody := autocompleteRequest{
		Input:                input,
		IncludedPrimaryTypes: []string{"(cities)"},
		SessionToken:         sessionToken,
		LanguageCode:         "en",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal autocomplete request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		googlePlacesBaseURL+"/v1/places:autocomplete",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create autocomplete request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("autocomplete request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"google autocomplete returned status %d",
			resp.StatusCode,
		)
	}

	var result autocompleteResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode autocomplete response: %w", err)
	}

	suggestions := make(
		[]models.LocationSuggestion,
		0,
		len(result.Suggestions),
	)

	for _, suggestion := range result.Suggestions {
		suggestions = append(
			suggestions,
			models.LocationSuggestion{
				PlaceID: suggestion.PlacePrediction.PlaceID,
				Text:    suggestion.PlacePrediction.Text.Text,
			},
		)
	}

	return suggestions, nil
}

func (c *GooglePlacesClient) GetPlaceDetails(
	placeID string,
	sessionToken string,
) (*models.Location, error) {
	query := url.Values{}
	query.Set("sessionToken", sessionToken)
	query.Set("languageCode", "en")

	requestURL := fmt.Sprintf(
		"%s/v1/places/%s?%s",
		googlePlacesBaseURL,
		url.PathEscape(placeID),
		query.Encode(),
	)

	req, err := http.NewRequest(
		http.MethodGet,
		requestURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create place details request: %w", err)
	}

	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", "addressComponents")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("place details request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"google place details returned status %d",
			resp.StatusCode,
		)
	}

	var result placeDetailsResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode place details response: %w", err)
	}

	location := &models.Location{}

	for _, component := range result.AddressComponents {
		for _, componentType := range component.Types {
			switch componentType {
			case "locality":
				location.City = component.LongText
			case "country":
				location.CountryCode = component.ShortText
			}
		}
	}

	if location.City == "" || location.CountryCode == "" {
		return nil, fmt.Errorf(
			"google place details missing city or country code",
		)
	}

	return location, nil
}
