package models

// APIKey represents an API Key object. The actual key value is located in the sys.id field.
type APIKey struct {
	System `json:"sys"`
	Name   string
}
