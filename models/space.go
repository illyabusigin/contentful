package models

import (
	"fmt"
)

// Space is a container for content types, entries and assets and other resources.
// API consumers, like mobile apps or websites, typically fetch data by getting
// entries and assets from one or more spaces.
type Space struct {
	System `json:"sys"`
	Name   string `json:"name"`
}

// Validate will validate the space. An error is returned if the space type is
// not  valid.
func (s *Space) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("Space must specify a valid name")
	}

	return nil
}

// Link returns a link to the space
func (s *Space) Link() *Link {
	return &Link{
		LinkData: &LinkData{
			Type:     LinkType,
			LinkType: "Space",
			ID:       s.ID,
		},
	}
}
