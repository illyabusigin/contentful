package management

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

// FetchAllLocales returns all locales associated with the provided space identifier
func (c *Client) FetchAllLocales(spaceID string) (locales []*Locale, pagination *Pagination, err error) {
	c.rl.Wait()

	type localesResponse struct {
		*Pagination
		Sys struct {
			Type string `json:"type"`
		} `json:"sys"`
		Items []*Locale `json:"items"`
	}

	results := new(localesResponse)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales", spaceID)
	_, err = c.sling.New().Get(path).Receive(results, contentfulError)

	return results.Items, results.Pagination, handleError(err, contentfulError)
}

// CreateLocale will create a locale with the provided information. It's important
// to note that you cannot create two lcoales with the same locale code.
func (c *Client) CreateLocale(spaceID string, locale *Locale) (created *Locale, err error) {
	if locale == nil {
		return nil, fmt.Errorf("CreateLocale failed, locale cannot be nil!")
	}
	if err = locale.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	// Default cannot be set via the API, set to false so it will not appear in the request body
	locale.Default = false

	created = new(Locale)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales", spaceID)
	_, err = c.sling.New().Post(path).BodyJSON(locale).Receive(created, contentfulError)

	return created, handleError(err, contentfulError)
}

// FetchLocale will return a locale for the given space and locale identifier.
func (c *Client) FetchLocale(spaceID string, localeID string) (locale *Locale, err error) {
	c.rl.Wait()

	locale = new(Locale)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales/%v", spaceID, localeID)
	_, err = c.sling.New().Get(path).Receive(locale, contentfulError)

	return locale, handleError(err, contentfulError)
}

// UpdateLocale will update the locale
func (c *Client) UpdateLocale(locale *Locale) (updated *Locale, err error) {
	if locale == nil {
		return nil, fmt.Errorf("Unable to locale. Locale argument was nil!")
	}

	// Default cannot be set via the API, set to false so it will not appear in the request body
	locale.Default = false

	if err = locale.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	updated = new(Locale)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales/%v", locale.Space.ID, locale.System.ID)
	_, err = c.sling.New().
		Set("X-Contentful-Version", fmt.Sprintf("%v", locale.System.Version)).
		Put(path).
		BodyJSON(locale).
		Receive(updated, contentfulError)

	return updated, handleError(err, contentfulError)
}

// DeleteLocale will delete an existing locale. Please note that it
// is not possible to recover from this action! Every content that
// was stored for that specific locale gets deleted and cannot be
// recreated by creating the same locale again.
func (c *Client) DeleteLocale(spaceID string, localeID string) (err error) {
	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/locales/%v", spaceID, localeID)
	_, err = c.sling.New().Delete(path).Receive(nil, contentfulError)

	return handleError(err, contentfulError)
}
