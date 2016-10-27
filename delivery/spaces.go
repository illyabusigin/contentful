package delivery

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// FetchSpace will return a space for the given identifier.
func (c *Client) FetchSpace(identifier string) (space *Space, err error) {
	c.rl.Wait()

	space = new(Space)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v", identifier)
	_, err = c.sling.New().Get(path).Receive(space, contentfulError)

	return space, handleError(err, contentfulError)
}
