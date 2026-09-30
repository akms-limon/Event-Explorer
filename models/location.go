package models

type LocationSuggestion struct {
	PlaceID string `json:"placeId"`
	Text    string `json:"text"`
}

type Location struct {
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
}
