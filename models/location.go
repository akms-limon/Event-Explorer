package models

type LocationSuggestion struct {
	PlaceID string
	Text    string
}

type Location struct {
	City        string
	CountryCode string
}