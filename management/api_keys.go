package management

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// CreateContentDeliveryAPIKey will create an API key with the specified name that can be used
// with the Content Delivery API (CDA).
func (c *Client) CreateContentDeliveryAPIKey(spaceID string, name string) (key *APIKey, err error) {
	if spaceID == "" {
		err = fmt.Errorf("CreateContentDeliveryAPIKey failed, spaceID cannot be empty!")
		return
	}

	if name == "" {
		err = fmt.Errorf("CreateContentDeliveryAPIKey failed, name cannot be empty!")
		return
	}

	c.rl.Wait()

	type apikey struct {
		name string
	}

	key = new(APIKey)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/api_keys", spaceID)
	_, err = c.sling.New().
		Post(path).
		BodyJSON(&apikey{name: name}).
		Receive(key, contentfulError)

	return key, handleError(err, contentfulError)
}

// FetchContentDeliveryAPIKeys returns all API Keys for the given space.
func (c *Client) FetchContentDeliveryAPIKeys(spaceID string, limit int, offset int) (keys []*APIKey, pagination *Pagination, err error) {
	if spaceID == "" {
		return nil, nil, fmt.Errorf("FetchContentDeliveryAPIKeys failed. Space identifier is not valid!")
	}

	if limit < 0 {
		return nil, nil, fmt.Errorf("FetchContentDeliveryAPIKeys failed. Limit must be greater than 0")
	}

	if limit > 100 {
		limit = 100
	}

	c.rl.Wait()

	type keysResponse struct {
		*Pagination
		Items []*APIKey `json:"items"`
	}

	results := new(keysResponse)
	results.Items = []*APIKey{}
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/api_keys", spaceID)
	req, err := c.sling.New().
		Get(path).
		Request()

	if err != nil {
		return
	}

	// Add query parameters
	q := req.URL.Query()

	q.Set("skip", fmt.Sprintf("%v", offset))
	q.Set("limit", fmt.Sprintf("%v", limit))
	req.URL.RawQuery = q.Encode()

	// Perform request
	_, err = c.sling.Do(req, results, contentfulError)

	return results.Items, results.Pagination, handleError(err, contentfulError)
}
