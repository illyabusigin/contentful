package management

import (
	"fmt"
	"time"
)

// Locale allow the definition of translated content for both assets and
// entries. A locale consists mainly of a name and a locale code.
type Locale struct {
	System   `json:"-"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Default  bool   `json:"default"`
	Optional bool   `json:"optional"`
	Fallback string `json:"fallbackCode"`

	SpaceID string `json:"-"`
}

func (l *Locale) Validate() error {
	if l.SpaceID == "" {
		return fmt.Errorf("Locale must specify valid SpaceID")
	}

	return nil
}

type contentfulLocale struct {
	Name                 string `json:"name"`
	InternalCode         string `json:"internal_code"`
	Code                 string `json:"code"`
	FallbackCode         string `json:"fallbackCode"`
	Default              bool   `json:"default"`
	ContentManagementAPI bool   `json:"contentManagementApi"`
	ContentDeliveryAPI   bool   `json:"contentDeliveryApi"`
	Optional             bool   `json:"optional"`
	Sys                  struct {
		Type    string `json:"type"`
		ID      string `json:"id"`
		Version int    `json:"version"`
		Space   struct {
			Sys struct {
				Type     string `json:"type"`
				LinkType string `json:"linkType"`
				ID       string `json:"id"`
			} `json:"sys"`
		} `json:"space"`
		CreatedBy struct {
			Sys struct {
				Type     string `json:"type"`
				LinkType string `json:"linkType"`
				ID       string `json:"id"`
			} `json:"sys"`
		} `json:"createdBy"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedBy struct {
			Sys struct {
				Type     string `json:"type"`
				LinkType string `json:"linkType"`
				ID       string `json:"id"`
			} `json:"sys"`
		} `json:"updatedBy"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"sys"`
}

func (c *contentfulLocale) Convert() *Locale {
	locale := new(Locale)

	locale.ID = c.Sys.ID
	locale.CreatedAt = c.Sys.CreatedAt
	locale.UpdatedAt = c.Sys.UpdatedAt

	locale.Type = c.Sys.Type
	locale.Version = c.Sys.Version

	locale.Name = c.Name
	locale.Code = c.Code
	locale.Default = c.Default
	locale.Optional = c.Optional
	locale.Fallback = c.FallbackCode
	locale.SpaceID = c.Sys.Space.Sys.ID

	return locale
}

type AllLocalesResponse struct {
	Locales []*Locale
	Error   error
}

// GetAllLocales returns all locales associated with the provided space identifier
func (c *Client) GetAllLocales(spaceIdentifier string) (response AllLocalesResponse) {
	c.rl.Wait()

	type localesResponse struct {
		Total int `json:"total"`
		Limit int `json:"limit"`
		Skip  int `json:"skip"`
		Sys   struct {
			Type string `json:"type"`
		} `json:"sys"`
		Items []contentfulLocale `json:"items"`
	}

	localesData := new(localesResponse)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales", spaceIdentifier)
	_, err := c.sling.New().Get(path).Receive(localesData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err != nil {
		response.Error = err
	} else {
		response.Locales = make([]*Locale, 0)

		for _, locale := range localesData.Items {
			response.Locales = append(response.Locales, locale.Convert())
		}
	}

	return
}

// CreateLocale will create a locale with the provided information. It's important
// to note that you cannot create two lcoales with the same locale code.
func (c *Client) CreateLocale(locale *Locale) (created *Locale, err error) {
	if err = locale.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	localeData := new(contentfulLocale)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales", locale.SpaceID)
	_, err = c.sling.New().Post(path).BodyJSON(locale).Receive(localeData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		created = localeData.Convert()
	}

	return
}

// FetchLocale will return a locale for the given space and locale identifier.
func (c *Client) FetchLocale(spaceIdentifier string, localeIdentifier string) (locale *Locale, err error) {
	c.rl.Wait()

	localeData := new(contentfulLocale)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales/%v", spaceIdentifier, localeIdentifier)
	_, err = c.sling.New().Get(path).Receive(localeData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		locale = localeData.Convert()
	}

	return
}

// UpdateLocale will update the locale
func (c *Client) UpdateLocale(locale *Locale) (updated *Locale, err error) {
	if locale == nil {
		return nil, fmt.Errorf("Unable to locale. Locale argument was nil!")
	}

	c.rl.Wait()

	localeData := new(contentfulLocale)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales/%v", locale.SpaceID, locale.System.ID)
	_, err = c.sling.New().
		Set("X-Contentful-Version", fmt.Sprintf("%v", locale.System.Version)).
		Put(path).
		BodyJSON(locale).
		Receive(localeData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		updated = localeData.Convert()
	}

	return
}

// DeleteLocale will delete an existing locale. Please note that it
// is not possible to recover from this action! Every content that
// was stored for that specific locale gets deleted and cannot be
// recreated by creating the same locale again.
func (c *Client) DeleteLocale(spaceIdentifier string, localIdentifier string) (err error) {
	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales/%v", spaceIdentifier, localIdentifier)
	_, err = c.sling.New().Delete(path).Receive(nil, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}
