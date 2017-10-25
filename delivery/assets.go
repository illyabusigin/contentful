package delivery

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// FetchAsset will return the specified asset.
func (c *Client) FetchAsset(spaceID string, assetID string) (asset *Asset, err error) {
	c.rl.Wait()

	asset = new(Asset)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v", spaceID, assetID)
	req, err := c.sling.New().
		Get(path).Request()

	// Add query parameters
	q := req.URL.Query()
	q.Set("locale", "*")

	req.URL.RawQuery = q.Encode()

	_, err = c.sling.Do(req, asset, contentfulError)

	return asset, handleError(err, contentfulError)
}

// QueryAssets will return all assets associated with a space
func (c *Client) QueryAssets(spaceID string, params map[string]string, limit int, offset int) (assets []*Asset, pagination *Pagination, err error) {
	if spaceID == "" {
		return nil, nil, fmt.Errorf("FetchAssets failed. Space identifier is not valid!")
	}

	if limit <= 0 {
		return nil, nil, fmt.Errorf("FetchAssets failed. Limit must be greater than 0")
	}

	if limit > PaginationSizeLimit {
		limit = PaginationSizeLimit
	}

	c.rl.Wait()

	type assetsResponse struct {
		*Pagination
		Items []*Asset `json:"items"`
	}

	results := new(assetsResponse)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets", spaceID)

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

	q.Set("locale", "*")
	q.Set("skip", fmt.Sprintf("%v", offset))
	q.Set("limit", fmt.Sprintf("%v", limit))
	req.URL.RawQuery = q.Encode()

	_, err = c.sling.Do(req, results, contentfulError)

	return results.Items, results.Pagination, handleError(err, contentfulError)
}
