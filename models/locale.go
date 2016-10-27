package models

import (
	"fmt"
)

// Locale allow the definition of translated content for both assets and
// entries. A locale consists mainly of a name and a locale code.
type Locale struct {
	System `json:"sys,omitempty"`

	Name     string `json:"name,omitempty"`
	Code     string `json:"code,omitempty"`
	Default  bool   `json:"default,omitempty"`
	Optional bool   `json:"optional"`
	Fallback string `json:"fallbackCode,omitempty"`

	EnabledForContentManagement bool `json:"contentManagementApi"`
	EnabledForContentDelivery   bool `json:"contentDeliveryApi"`
}

// Validate will validate the locale. An error is returned if the content type
// is not  valid.
func (l *Locale) Validate() error {
	if len(l.Name) == 0 {
		return fmt.Errorf("Locale name cannot be empty")
	}

	if len(l.Code) == 0 {
		return fmt.Errorf("Locale code cannot be empty")
	}

	return nil
}
