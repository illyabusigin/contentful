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
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	LinkType string    `json:"linkType,omitempty"`
	Type     FieldType `json:"type,omitempty"`

	Items *Field `json:"items,omitempty"`

	Localized   bool              `json:"localized,omitempty"`
	Required    bool              `json:"required,omitempty"`
	Validations []FieldValidation `json:"validations,omitempty"`
	Disabled    bool              `json:"disabled,omitempty"`

	// Omitted fields will stil be present in CMA APIs but omitted from CDA and CPA APIs
	Omitted bool `json:"omitted,omitempty"`
}

// FieldValidation describes validation rules associated with a field, if any.
type FieldValidation struct {
	Size                 *SizeFieldValidation      `json:"size,omitempty"`
	Range                *RangeFieldValidation     `json:"range,omitempty"`
	DateRange            *DateRangeFieldValidation `json:"dateRange,omitempty"`
	RegularExpression    *RegExFieldValidation     `json:"regexp,omitempty"`
	AssetImageValidation *AssetImageValidation     `json:"assetImageDimensions,omitempty"`

	LinkMIMETypeGroup []string      `json:"linkMimetypeGroup,omitempty"`
	LinkContentTypes  []string      `json:"linkContentType,omitempty"`
	In                []interface{} `json:"in,omitempty"`
	Message           string        `json:"message,omitempty"`
}

// RegExFieldValidation permits validation with regular expression.
type RegExFieldValidation struct {
	Pattern string `json:"pattern"`
}

// RangeFieldValidation takes optional min and max parameters and validates the range
// of a value.
type RangeFieldValidation struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

// SizeFieldValidation permits validation with size. You can specify
// either minimum size, maximum size, or both.
type SizeFieldValidation struct {
	Min float64 `json:"min,omitempty"`
	Max float64 `json:"max,omitempty"`
}

// AssetImageValidation permits asset validation around size
type AssetImageValidation struct {
	Width  *SizeFieldValidation `json:"width,omitempty"`
	Height *SizeFieldValidation `json:"height,omitempty"`
}

// DateRangeFieldValidation permits validation with date ranges. You can specify
// either minimum date, maximum date, or both.
type DateRangeFieldValidation struct {
	Min *time.Time `json:"min,omitempty"`
	Max *time.Time `json:"max,omitempty"`
}

// ContentType are schemas that define the fields of entries.
// Every entry can only contain values in the fields defined by
// its content type, and the values of those fields must match
// the data type defined in the content type. There is a limit
// of 50 fields per content type.
type ContentType struct {
	System `json:"sys"`

	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	DisplayField string  `json:"displayField,omitempty"`
	Fields       []Field `json:"fields"`
}

// Validate will validate the content type. An error is returned if the content
// type is not  valid.
func (t *ContentType) Validate() error {
	if len(t.Name) == 0 {
		return fmt.Errorf("Content type name cannot be empty")
	}

	if t.Space == nil || t.Space.ID == "" {
		return fmt.Errorf("Locale must specify valid System.Space.ID!")
	}

	if len(t.ID) == 0 {
		return fmt.Errorf("Content type must specify an identifier!")
	}

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

	req, err := c.sling.New().
		Get(path()).
		Request()

	if err != nil {
		return
	}

	// Add query parameters
	q := req.URL.Query()
	q.Set("skip", fmt.Sprintf("%v", offset))
	q.Set("limit", fmt.Sprintf("%v", limit))
	req.URL.RawQuery = q.Encode()

	_, err = c.sling.Do(req, results, contentfulError)

	return results.Items, results.Pagination, handleError(err, contentfulError)

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

	created = new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", contentType.Space.ID, contentType.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		BodyJSON(contentType).
		Receive(created, contentfulError)

	return created, handleError(err, contentfulError)
}

// FetchContentType will return a content type for the specified space and
// content type identifier.
func (c *Client) FetchContentType(spaceID string, contentTypeID string) (contentType *ContentType, err error) {
	if spaceID == "" || contentTypeID == "" {
		err = fmt.Errorf("FetchContentType failed. Invalid spaceID or contentTypeID.")
		return
	}

	c.rl.Wait()

	contentType = new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", spaceID, contentTypeID)
	_, err = c.sling.New().Get(path).Receive(contentType, contentfulError)

	return contentType, handleError(err, contentfulError)
}

// UpdateContentType will update the content type with the specified changes.
func (c *Client) UpdateContentType(contentType *ContentType) (updated *ContentType, err error) {
	return c.CreateContentType(contentType)
}

// DeleteContentType will delete a content type. Before you can delete a content
// type you need to deactivate it.
func (c *Client) DeleteContentType(spaceID string, contentTypeID string) (err error) {
	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", spaceID, contentTypeID)
	_, err = c.sling.New().
		Delete(path).
		Receive(nil, contentfulError)

	return handleError(err, contentfulError)
}

// ActivateContentType makes the content type available for creating entries
func (c *Client) ActivateContentType(contentType *ContentType) (activated *ContentType, err error) {
	if contentType == nil {
		return nil, fmt.Errorf("CreateContentType failed. Type argument was nil!")
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	activated = new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v/published", contentType.Space.ID, contentType.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		Receive(activated, contentfulError)

	return activated, handleError(err, contentfulError)
}

// DeactivateContentType removes the availability for creating entries
func (c *Client) DeactivateContentType(contentType *ContentType) (deactivated *ContentType, err error) {
	if contentType == nil {
		return nil, fmt.Errorf("CreateContentType failed. Type argument was nil!")
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	deactivated = new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v/published", contentType.Space.ID, contentType.ID)
	_, err = c.sling.New().
		Delete(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		Receive(deactivated, contentfulError)

	return deactivated, handleError(err, contentfulError)
}
