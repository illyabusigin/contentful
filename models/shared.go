package models

import (
	"time"
)

// System contains system managed metadata. The exact metadata available depends
// on the type of the resource but at minimum System.Type property is defined.
// Note that none of the sys fields are editable and only the sys.id field can
// be specified in the creation of an item (as long as it is not a Space).
type System struct {
	ID        string     `json:"id,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Type    string `json:"type,omitempty"`
	Version int    `json:"version,omitempty"`

	Space       *Link `json:"space,omitempty"`
	ContentType *Link `json:"contentType,omitempty"`

	FirstPublished   *time.Time `json:"firstPublishedAt,omitempty"`
	PublishedAt      *time.Time `json:"publishedAt,omitempty"`
	PublishedVersion int        `json:"publishedVersion,omitempty"`
	ArchivedAt       *time.Time `json:"archivedAt,omitempty"`
}

// Link represents a link to another Contentful object
type Link struct {
	*LinkData `json:"sys"`
}

// LinkData contains the link information
type LinkData struct {
	Type     string `json:"type"`
	LinkType string `json:"linkType"`
	ID       string `json:"id"`
}

// Pagination represents all paginated data and is returned for collection
// related client methods.
type Pagination struct {
	Total int
	Skip  int
	Limit int
}
