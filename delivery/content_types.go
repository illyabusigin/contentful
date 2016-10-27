package delivery

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// FetchContentTypes returns all content types for a given space. You can filter
// this further by toggling the published flag
func (c *Client) FetchContentTypes(spaceID string, limit int, offset int) (contentTypes []*ContentType, pagination *Pagination, err error) {
	if spaceID == "" {
		return nil, nil, fmt.Errorf("FetchContentTypes failed. Space identifier is not valid!")
	}

	if limit <= 0 {
		return nil, nil, fmt.Errorf("FetchContentTypes failed. Limit must be greater than 0")
	}

	if limit > 100 {
		limit = 100
	}

	c.rl.Wait()

	type contentTypesResponse struct {
		*Pagination
		Items []*ContentType `json:"items"`
	}

	results := new(contentTypesResponse)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types", spaceID)

	req, err := c.sling.New().
		Get(path).Request()

	if err != nil {
		return
	}

	// Add query parameters
	q := req.URL.Query()
	q.Set("skip", fmt.Sprintf("%v", offset))
	q.Set("limit", fmt.Sprintf("%v", limit))
	req.URL.RawQuery = q.Encode()

	_, err = c.sling.Do(req, results, contentfulError)

	return results.Items, results.Pagination, handleError(err, contentfulError)
}

// FetchContentType will return a content type for the specified space and
// content type identifier.
func (c *Client) FetchContentType(spaceID string, contentTypeID string) (contentType *ContentType, err error) {
	if spaceID == "" || contentTypeID == "" {
		err = fmt.Errorf("FetchContentType failed. Invalid spaceID or contentTypeID.")
		return
	}

	c.rl.Wait()

	contentType = new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", spaceID, contentTypeID)
	_, err = c.sling.New().Get(path).Receive(contentType, contentfulError)

	return contentType, handleError(err, contentfulError)
}
