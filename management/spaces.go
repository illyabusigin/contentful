package management

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

func (s *Space) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("Space must specify a valid name")
	}

	return nil
}

// GetAllSpaces returns all spaces associated with the account
func (c *Client) FetchAllSpaces() (spaces []*Space, pagination *Pagination, err error) {
	c.rl.Wait()

	type spacesResponse struct {
		*Pagination
		Sys struct {
			Type string `json:"type"`
		} `json:"sys"`
		Items []*Space `json:"items"`
	}

	results := new(spacesResponse)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces")
	_, err = c.sling.New().Get(path).Receive(results, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err != nil {
		return nil, nil, err
	}

	return spaces, results.Pagination, nil
}

// CreateSpace will create a space with the provided name. It's important to
// note that names are not unique between spaces.
func (c *Client) CreateSpace(space *Space) (created *Space, err error) {
	if err = space.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	created = new(Space)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces")
	_, err = c.sling.New().Post(path).BodyJSON(space).Receive(created, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// FetchSpace will return a space for the given identifier.
func (c *Client) FetchSpace(identifier string) (space *Space, err error) {
	c.rl.Wait()

	space = new(Space)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v", identifier)
	_, err = c.sling.New().Get(path).Receive(space, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// UpdateSpace will update the space
func (c *Client) UpdateSpace(space *Space) (updated *Space, err error) {
	if space == nil {
		return nil, fmt.Errorf("Unable to update. Space argument was nil!")
	}

	c.rl.Wait()

	updated = new(Space)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v", space.System.ID)
	_, err = c.sling.New().
		Set("X-Contentful-Version", fmt.Sprintf("%v", space.System.Version)).
		Put(path).
		BodyJSON(space).
		Receive(updated, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// DeleteSpace will delete an existing space by doing a DELETE request to /spaces/ID.
// Note that deleting a space will remove its entire content, including all content
// types, entries and assets. Be careful as this action can not be undone.
func (c *Client) DeleteSpace(identifier string) (err error) {
	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v", identifier)
	_, err = c.sling.New().Delete(path).Receive(nil, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}
