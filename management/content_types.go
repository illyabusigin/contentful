package management

import (
	"fmt"
	"time"
)

// Content Field Types
type FieldType string

// Content Field Type constants
const (
	// ShortText fields do not support ordering or strict equality. 1 to 256 characters.
	ShortText FieldType = "Symbol"
	// LongText can be a maximum length of 50,000 characters
	LongText = "Text"
	Integer  = "Integer"
	Number   = "Number"
	Boolean  = "Boolean"
	// Date must be ISO8601 formatted and do not require a time portion
	Date     = "Date"
	Location = "Location"
	Object   = "Object"
	LinkType = "Link"
	Array    = "Array"
)

// Field describes a single allowed field value of an an entry.
// Each field type corresponds to a JSON type, though there are
// more field types than JSON types.
type Field struct {
	ID          string
	Name        string
	Type        FieldType
	Localized   bool
	Required    bool
	Validations []FieldValidation
	Disabled    bool

	// Omitted fields will stil be present in CMA APIs but omitted from CDA and CPA APIs
	Omitted bool
}

// FieldValidation describes validation rules associated with a field, if any.
type FieldValidation struct {
	Size *struct {
		Min *float64
		Max *float64
	} `json:"size,omitempty"`

	DateRange *struct {
		Min *time.Time
		Max *time.Time
	} `json:"dateRange,omitempty"`

	RegularExpression *struct {
		Pattern string
	} `json:"regexp,omitempty"`

	LinkMIMETypeGroup string        `json:"linkMimetypeGroup,omitempty"`
	LinkContentTypes  []string      `json:"linkContentType,omitempty"`
	In                []interface{} `json:"in,omitempty"`
	Message           *string       `json:"message,omitempty"`
}

// ContentType are schemas that define the fields of entries.
// Every entry can only contain values in the fields defined by
// its content type, and the values of those fields must match
// the data type defined in the content type. There is a limit
// of 50 fields per content type.
type ContentType struct {
	System  `json:"sys"`
	SpaceID string `json:"-"`

	Name         string  `json:"-"`
	Description  string  `json:"description,omitempty"`
	DisplayField string  `json:"displayField,omitempty"`
	Fields       []Field `json:"fields"`
}

// Validate the ContentType
func (t *ContentType) Validate() error {
	return nil
}

// FetchContentTypes returns all content types for a given space. You can filter
// this further by toggling the published flag
func (c *Client) FetchContentTypes(spaceID string, published bool, limit int, offset int) (contentTypes []*ContentType, pagination *Pagination, err error) {
	if spaceID == "" {
		return nil, nil, fmt.Errorf("FetchContentTypes failed. Space identifier is not valid!")
	}

	if limit <= 0 {
		return nil, nil, fmt.Errorf("FetchContentTypes failed. Limit must be greater than 0")
	}

	if limit > 100 {
		limit = 100
	}

	c.rl.Wait()

	type contentTypesResponse struct {
		*Pagination
		Items []*ContentType `json:"items"`
	}

	results := new(contentTypesResponse)
	contentfulError := new(ContentfulError)
	path := func() string {
		if published {
			return fmt.Sprintf("spaces/%v/public/content_types", spaceID)
		}

		return fmt.Sprintf("spaces/%v/content_types", spaceID)
	}
	_, err = c.sling.New().
		Get(path()).
		Receive(results, contentfulError)

	if err != nil {
		return
	}

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	for _, contentType := range results.Items {
		contentType.SpaceID = spaceID
	}

	return results.Items, results.Pagination, nil

}

// CreateContentType will create a content type. It's recommended that you
// control the ID of the created content type and associated fields.
func (c *Client) CreateContentType(contentType *ContentType) (created *ContentType, err error) {
	if contentType == nil {
		return nil, fmt.Errorf("CreateContentType failed. Type argument was nil!")
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	createdData := new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", contentType.SpaceID, contentType.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		BodyJSON(contentType).
		Receive(createdData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		created = createdData
		created.SpaceID = contentType.SpaceID
	}

	return
}

// FetchContentType will return a content type for the specified space and
// content type identifier.
func (c *Client) FetchContentType(spaceID string, contentTypeID string) (contentType *ContentType, err error) {
	if spaceID == "" || contentTypeID == "" {
		err = fmt.Errorf("FetchContentType failed. Invalid spaceID or contentTypeID.")
		return
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", spaceID, contentTypeID)
	_, err = c.sling.New().Get(path).Receive(contentType, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// UpdateContentType will update the content type with the specified changes.
func (c *Client) UpdateContentType(contentType *ContentType) (updated *ContentType, err error) {
	return c.CreateContentType(contentType)
}

// DeleteContentType will delete a content type. Before you can delete a content
// type you need to deactivate it.
func (c *Client) DeleteContentType(contentType *ContentType) (err error) {
	if contentType == nil {
		return fmt.Errorf("DeleteContentType failed. Type argument was nil!")
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", contentType.SpaceID, contentType.ID)
	_, err = c.sling.New().
		Delete(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		BodyJSON(contentType).
		Receive(nil, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// ActivateContentType makes the content type available for creating entries
func (c *Client) ActivateContentType(contentType *ContentType) (current *ContentType, err error) {
	if contentType == nil {
		return nil, fmt.Errorf("CreateContentType failed. Type argument was nil!")
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	publishedType := new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v/published", contentType.SpaceID, contentType.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		Receive(publishedType, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		current = publishedType
		current.SpaceID = publishedType.SpaceID
	}

	return
}

// DeactivateContentType removes the availability for creating entries
func (c *Client) DeactivateContentType(contentType *ContentType) (current *ContentType, err error) {
	if contentType == nil {
		return nil, fmt.Errorf("CreateContentType failed. Type argument was nil!")
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	publishedType := new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v/published", contentType.SpaceID, contentType.ID)
	_, err = c.sling.New().
		Delete(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		Receive(publishedType, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		current = publishedType
		current.SpaceID = publishedType.SpaceID
	}

	return
}
