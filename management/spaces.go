package management

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// FetchAllSpaces returns all spaces associated with the account
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

	return spaces, results.Pagination, handleError(err, contentfulError)
}

// CreateSpace will create a space with the provided name. It's important to
// note that names are not unique between spaces.
func (c *Client) CreateSpace(space *Space) (created *Space, err error) {
	if space == nil {
		return nil, fmt.Errorf("CreateSpace failed. Space cannot be nil!")
	}

	if err = space.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	created = new(Space)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces")
	_, err = c.sling.New().Post(path).BodyJSON(space).Receive(created, contentfulError)

	return created, handleError(err, contentfulError)
}

// FetchSpace will return a space for the given identifier.
func (c *Client) FetchSpace(identifier string) (space *Space, err error) {
	c.rl.Wait()

	space = new(Space)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v", identifier)
	_, err = c.sling.New().Get(path).Receive(space, contentfulError)

	return space, handleError(err, contentfulError)
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

	return updated, handleError(err, contentfulError)
}

// DeleteSpace will delete an existing space by doing a DELETE request to /spaces/ID.
// Note that deleting a space will remove its entire content, including all content
// types, entries and assets. Be careful as this action can not be undone.
func (c *Client) DeleteSpace(identifier string) (err error) {
	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v", identifier)
	_, err = c.sling.New().Delete(path).Receive(nil, contentfulError)

	return handleError(err, contentfulError)
}
