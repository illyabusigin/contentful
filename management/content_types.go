package management

import (
	"fmt"
	"time"
)

// Content Field Types
type FieldType string

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
	Link     = "Link"
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
	Omitted     bool
}

// FieldValidation describes the validation rules associated with
// a field, if any.
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

func (t *ContentType) Validate() error {
	return nil
}

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
