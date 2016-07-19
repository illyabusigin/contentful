package management

import (
	"fmt"
	"time"
)

// Space is a container for content types, entries and assets and other resources.
// API consumers, like mobile apps or websites, typically fetch data by getting
// entries and assets from one or more spaces.
type Space struct {
	System `json:"-"`
	Name   string `json:"name"`
}

func (s *Space) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("Space must specify a valid name")
	}

	return nil
}

type contentfulSpace struct {
	Name string `json:"name"`
	Sys  struct {
		Type      string `json:"type"`
		ID        string `json:"id"`
		Version   int    `json:"version"`
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

func (c *contentfulSpace) Convert() *Space {
	space := new(Space)

	space.ID = c.Sys.ID
	space.CreatedAt = c.Sys.CreatedAt
	space.UpdatedAt = c.Sys.UpdatedAt

	space.Type = c.Sys.Type
	space.Version = c.Sys.Version

	space.Name = c.Name

	return space
}

type AllSpacesResponse struct {
	Spaces []*Space
	Error  error
}

// GetAllSpaces returns all spaces associated with the account
func (c *Client) GetAllSpaces() (response AllSpacesResponse) {
	c.rl.Wait()

	type SpacesResponse struct {
		Total int `json:"total"`
		Limit int `json:"limit"`
		Skip  int `json:"skip"`
		Sys   struct {
			Type string `json:"type"`
		} `json:"sys"`
		Items []contentfulSpace `json:"items"`
	}

	spacesData := new(SpacesResponse)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces")
	_, err := c.sling.New().Get(path).Receive(spacesData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err != nil {
		response.Error = err
	} else {
		response.Spaces = make([]*Space, 0)

		for _, space := range spacesData.Items {
			response.Spaces = append(response.Spaces, space.Convert())
		}
	}

	return
}

// CreateSpace will create a space with the provided name. It's important to
// note that names are not unique between spaces.
func (c *Client) CreateSpace(space *Space) (created *Space, err error) {
	if err = space.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	spaceData := new(contentfulSpace)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces")
	_, err = c.sling.New().Post(path).BodyJSON(space).Receive(spaceData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		created = spaceData.Convert()
	}

	return
}

// FetchSpace will return a space for the given identifier.
func (c *Client) FetchSpace(identifier string) (space *Space, err error) {
	c.rl.Wait()

	spaceData := new(contentfulSpace)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v", identifier)
	_, err = c.sling.New().Get(path).Receive(spaceData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		space = spaceData.Convert()
	}

	return
}

// UpdateSpace will update the space
func (c *Client) UpdateSpace(space *Space) (updated *Space, err error) {
	if space == nil {
		return nil, fmt.Errorf("Unable to update. Space argument was nil!")
	}

	c.rl.Wait()

	spaceData := new(contentfulSpace)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v", space.System.ID)
	_, err = c.sling.New().
		Set("X-Contentful-Version", fmt.Sprintf("%v", space.System.Version)).
		Put(path).
		BodyJSON(space).
		Receive(spaceData, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	if err == nil {
		updated = spaceData.Convert()
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
