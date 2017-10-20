package models

import "fmt"

// Error represnts the error object that is returned when something
// goes wrong with a Contentful API request. This struct conforms to the `error`
// interface.
type Error struct {
	error
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
	Sys       struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"sys"`
	Details struct {
		Errors []interface{} `json:"errors"`
	} `json:"details"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v, %v, %v", e.Message, e.RequestID, e.Sys)
}
