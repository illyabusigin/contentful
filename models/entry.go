package models

import (
	"fmt"
)

type Linkable interface {
	Link() map[string]map[string]interface{}
}

type EntryFields map[string]interface{}

type NewEntry struct {
	Fields EntryFields `json:"fields"`
}

// Validate validates the entry
func (c *NewEntry) Validate() error {
	if c.Fields == nil || len(c.Fields) == 0 {
		return fmt.Errorf("NewEntry.Fields cannot be empty!")
	}

	return nil
}

// Includes are a sub-type used by EntryResults
type Includes struct {
	Entries []*Entry `json:"Entry"`
	Assets  []*Asset `json:"Asset"`
}

// QueryEntriesResult are returned for QueryEntries
type QueryEntriesResult struct {
	Entries  []*Entry
	Includes *Includes

	Pagination *Pagination
	Errors     []error
}

// Entry represent textual content in a space. An entry's data adheres to a
// certain content type.
type Entry struct {
	System `json:"sys"`
	Fields EntryFields `json:"fields"`
}

// Validate will validate the entry. An error is returned if the entry
// is not  valid.
func (c *Entry) Validate() error {
	if c.Space == nil || c.Space.ID == "" {
		return fmt.Errorf("Entry must have a valid Space associated with it!")
	}

	if c.System.ID == "" {
		return fmt.Errorf("Entry.System.ID cannot be empty!")
	}

	if c.Fields == nil || len(c.Fields) == 0 {
		return fmt.Errorf("Entry.Fields cannot be empty!")
	}

	return nil
}

// Link represents a link to the Entry and implements the Linkable interface
func (c *Entry) Link() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"sys": map[string]interface{}{
			"id":       c.ID,
			"linkType": "Entry",
			"type":     LinkType,
		},
	}
}
