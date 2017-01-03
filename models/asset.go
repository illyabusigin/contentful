package models

import (
	"fmt"
)

// Asset represents a file within a space. An asset can be any kind of file: an
// image, a video, an audio file, a PDF or any other filetype. Assets are
// usually attached to entries through links.
//
// Assets can optionally be localized by providing separate files for each
// locale. Those assets which are not localized simply provide a single file
// under the default locale.
type Asset struct {
	System `json:"sys"`
	Fields AssetFields `json:"fields,omitempty"`
}

// AssetFields contains all asset information.
type AssetFields struct {
	Title map[string]string    `json:"title,omitempty"`
	File  map[string]AssetData `json:"file,omitempty"`
}

// AssetData contains all asset information
type AssetData struct {
	MIMEType string      `json:"contentType"`
	Name     string      `json:"fileName"`
	URL      string      `json:"url,omitempty"`
	Upload   string      `json:"upload,omitempty"`
	Detail   AssetDetail `json:"details,omitempty"`
}

// Validate will validate the Asset to ensure all necessary fields are present.
func (a *Asset) Validate() error {
	if a.System.ID == "" {
		fmt.Println("system ID is empty!")
		return fmt.Errorf("Asset validation failed. System.ID cannot be empty!")
	}

	if a.System.Space == nil {
		return fmt.Errorf("Asset validation failed. System.Space.ID cannot be empty!")
	}

	return nil
}

// Link represents a link to the Entry and implements the Linkable interface
func (a *Asset) Link() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"sys": map[string]interface{}{
			"id":       a.ID,
			"linkType": "Asset",
			"type":     LinkType,
		},
	}
}

type ImageDetail struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AssetDetail struct {
	Image ImageDetail `json:"image,omitempty"`
	Size  int         `json:"size"`
}

// File represents all asset data prior to upload
type File struct {
	SpaceID string `json:"-"`

	Fields FileFields `json:"fields"`
}

type FileFields struct {
	Title map[string]string   `json:"title"`
	File  map[string]FileData `json:"file"`
}

// FileData contains all file information
type FileData struct {
	MIMEType string `json:"contentType"`
	Name     string `json:"fileName,omitempty"`
	URL      string `json:"upload,omitempty"`
}

func (f *File) Validate() error {
	if f.SpaceID == "" {
		return fmt.Errorf("Filed validation failed. SpaceID cannot be empty!")
	}

	if f.Fields.File == nil || len(f.Fields.File) == 0 {
		return fmt.Errorf("Filed validation failed. Fields.File cannot be empty!")
	}

	if f.Fields.Title == nil || len(f.Fields.Title) == 0 {
		return fmt.Errorf("Filed validation failed. Fields.Title cannot be empty!")
	}

	for _, data := range f.Fields.File {
		if data.Name == "" {
			return fmt.Errorf("Filed validation failed. FileData.Name cannot be empty. FileData: %v", data)
		} else if data.MIMEType == "" {
			return fmt.Errorf("Filed validation failed. FileData.MIMEType cannot be empty. FileData: %v", data)
		} else if data.URL == "" {
			return fmt.Errorf("Filed validation failed. FileData.URL cannot be empty. FileData: %v", data)
		}
	}

	return nil
}
