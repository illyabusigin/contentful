package delivery

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// QueryEntries returns all entries for the given space and parameters.
func (c *Client) QueryEntries(spaceID string, params map[string]string, limit int, offset int) (entries []*Entry, pagination *Pagination, err error) {
	if spaceID == "" {
		return nil, nil, fmt.Errorf("FetchEntries failed. Space identifier is not valid!")
	}

	if limit < 0 {
		return nil, nil, fmt.Errorf("FetchEntries failed. Limit must be greater than 0")
	}

	if limit > 100 {
		limit = 100
	}

	c.rl.Wait()

	type entriesResponse struct {
		*Pagination
		Items []*Entry `json:"items"`
	}

	results := new(entriesResponse)
	results.Items = []*Entry{}
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries", spaceID)
	req, err := c.sling.New().
		Get(path).
		Request()

	if err != nil {
		return
	}

	// Add query parameters
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}

	q.Set("skip", fmt.Sprintf("%v", offset))
	q.Set("limit", fmt.Sprintf("%v", limit))
	req.URL.RawQuery = q.Encode()

	// Perform request
	_, err = c.sling.Do(req, results, contentfulError)

	return results.Items, results.Pagination, handleError(err, contentfulError)
}

// FetchEntry returns a single entry for the given space and entry identifier
func (c *Client) FetchEntry(spaceID string, entryID string) (entry *Entry, err error) {
	if spaceID == "" || entryID == "" {
		err = fmt.Errorf("FetchContentType failed. Invalid spaceID or contentTypeID.")
		return
	}

	c.rl.Wait()

	entry = new(Entry)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries/%v", spaceID, entryID)
	_, err = c.sling.New().Get(path).Receive(entry, contentfulError)

	return entry, handleError(err, contentfulError)
}
