package delivery

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// QueryEntries returns all entries for the given space and parameters.
func (c *Client) QueryEntries(spaceID string, params map[string]string, limit int, offset int) (result *QueryEntriesResult) {
	result = &QueryEntriesResult{
		Entries: []*Entry{},
		Includes: &Includes{
			Entries: []*Entry{},
			Assets:  []*Asset{},
		},

		Errors: []error{},
	}

	if spaceID == "" {
		result.Errors = append(result.Errors, fmt.Errorf("QueryEntries failed. Space identifier is not valid!"))
		return
	}

	if limit < 0 {
		result.Errors = append(result.Errors, fmt.Errorf("QueryEntries failed. Limit must be greater than 0"))
		return
	}

	if limit > 100 {
		limit = 100
	}

	c.rl.Wait()

	type entriesResponse struct {
		*Pagination
		Items    []*Entry  `json:"items"`
		Includes *Includes `json:"includes"`
	}

	response := new(entriesResponse)
	response.Items = []*Entry{}
	response.Includes = &Includes{
		Entries: []*Entry{},
		Assets:  []*Asset{},
	}

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
	_, err = c.sling.Do(req, response, contentfulError)

	result.Pagination = response.Pagination
	result.Includes = response.Includes
	result.Entries = response.Items

	if handledErr := handleError(err, contentfulError); handledErr != nil {
		result.Errors = append(result.Errors, handledErr)
	}

	return
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
