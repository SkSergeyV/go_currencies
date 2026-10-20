package importer

import "time"

type JsonImport struct {
	Date         time.Time           `json:"Date"`
	PreviousDate time.Time           `json:"PreviousDate"`
	PreviousURL  string              `json:"PreviousURL"`
	Timestamp    time.Time           `json:"Timestamp"`
	Valute       map[string]Currency `json:"Valute"`
}

type Currency struct {
	ID       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Nominal  int     `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
	Previous float64 `json:"Previous"`
}
