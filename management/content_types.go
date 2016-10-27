package management

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// FetchContentTypes returns all content types for a given space. You can filter
// this further by toggling the published flag
func (c *Client) FetchContentTypes(spaceID string, published bool, limit int, offset int) (contentTypes []*ContentType, pagination *Pagination, err error) {
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
	path := func() string {
		if published {
			return fmt.Sprintf("spaces/%v/public/content_types", spaceID)
		}

		return fmt.Sprintf("spaces/%v/content_types", spaceID)
	}

	req, err := c.sling.New().
		Get(path()).
		Request()

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

// CreateContentType will create a content type. It's recommended that you
// control the ID of the created content type and associated fields.
func (c *Client) CreateContentType(contentType *ContentType) (created *ContentType, err error) {
	if contentType == nil {
		return nil, fmt.Errorf("CreateContentType failed. Type argument was nil!")
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	created = new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", contentType.Space.ID, contentType.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		BodyJSON(contentType).
		Receive(created, contentfulError)

	fmt.Println("created:", created)
	fmt.Println("err", err)
	fmt.Println("contentFulError", contentfulError)

	return created, handleError(err, contentfulError)
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

// UpdateContentType will update the content type with the specified changes.
func (c *Client) UpdateContentType(contentType *ContentType) (updated *ContentType, err error) {
	return c.CreateContentType(contentType)
}

// DeleteContentType will delete a content type. Before you can delete a content
// type you need to deactivate it.
func (c *Client) DeleteContentType(spaceID string, contentTypeID string) (err error) {
	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v", spaceID, contentTypeID)
	_, err = c.sling.New().
		Delete(path).
		Receive(nil, contentfulError)

	return handleError(err, contentfulError)
}

// ActivateContentType makes the content type available for creating entries
func (c *Client) ActivateContentType(contentType *ContentType) (activated *ContentType, err error) {
	if contentType == nil {
		return nil, fmt.Errorf("CreateContentType failed. Type argument was nil!")
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	activated = new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v/published", contentType.Space.ID, contentType.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		Receive(activated, contentfulError)

	return activated, handleError(err, contentfulError)
}

// DeactivateContentType removes the availability for creating entries
func (c *Client) DeactivateContentType(contentType *ContentType) (deactivated *ContentType, err error) {
	if contentType == nil {
		return nil, fmt.Errorf("CreateContentType failed. Type argument was nil!")
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	deactivated = new(ContentType)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/content_types/%v/published", contentType.Space.ID, contentType.ID)
	_, err = c.sling.New().
		Delete(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", contentType.Version)).
		Receive(deactivated, contentfulError)

	return deactivated, handleError(err, contentfulError)
}
