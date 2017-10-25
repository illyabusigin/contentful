package management

import (
	"fmt"

	. "github.com/illyabusigin/contentful/models"
)

// CreateAsset creates a new asset. It's important to note that the asset still
// needs to be processed and published to be availably through the delivery API.
func (c *Client) CreateAsset(file *File) (created *Asset, err error) {
	if file == nil {
		err = fmt.Errorf("CreateAsset failed. Entry must not be nil!")
		return
	}

	if err = file.Validate(); err != nil {
		return
	}

	c.rl.Wait()
	created = &Asset{}

	created = new(Asset)
	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets", file.SpaceID)
	_, err = c.sling.New().
		Post(path).
		BodyJSON(file).
		Receive(created, contentfulError)

	return created, handleError(err, contentfulError)
}

// UpdateAsset will update the asset
func (c *Client) UpdateAsset(asset *Asset) (updated *Asset, err error) {
	if asset == nil {
		err = fmt.Errorf("UpdateAsset failed. Asset must not be nil!")
		return
	}

	if err = asset.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	updated = new(Asset)
	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets/%v", asset.System.Space.ID, asset.System.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", asset.System.Version)).
		BodyJSON(asset).
		Receive(updated, contentfulError)

	return updated, handleError(err, contentfulError)
}

// FetchAsset will return the specified asset.
func (c *Client) FetchAsset(spaceID string, assetID string) (asset *Asset, err error) {
	c.rl.Wait()

	asset = new(Asset)
	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets/%v", spaceID, assetID)
	_, err = c.sling.New().
		Get(path).
		BodyJSON(asset).
		Receive(asset, contentfulError)

	return asset, handleError(err, contentfulError)
}

// QueryAssets will return all assets associated with a space. You can toggle
// the published flag to only fetch published assets.
func (c *Client) QueryAssets(spaceID string, published bool, params map[string]string, limit int, offset int) (assets []*Asset, pagination *Pagination, err error) {
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
	contentfulError := new(Error)
	path := func() string {
		if published {
			return fmt.Sprintf("spaces/%v/public/assets", spaceID)
		}

		return fmt.Sprintf("spaces/%v/assets", spaceID)
	}

	req, err := c.sling.New().
		Get(path()).
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

	_, err = c.sling.Do(req, results, contentfulError)

	return results.Items, results.Pagination, handleError(err, contentfulError)
}

// ProcessAsset process the asset. This uploads the asset to
// Contentful among other things. Processing happens asynchronously, the call
// will not block until it has finished.
func (c *Client) ProcessAsset(asset *Asset, localeCode string) (err error) {
	if asset == nil {
		return fmt.Errorf("ProcessAsset failed. Asset cannot be nil!")
	}

	if localeCode == "" {
		return fmt.Errorf("ProcessAsset failed. Locale cannot be empty!")
	}

	c.rl.Wait()

	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets/%v/files/%v/process", asset.Space.ID, asset.ID, localeCode)
	_, err = c.sling.New().
		Put(path).
		Receive(nil, contentfulError)

	return handleError(err, contentfulError)
}

// PublishAsset will make the asset available via the Content Delivery API.
func (c *Client) PublishAsset(asset *Asset) (published *Asset, err error) {
	if asset == nil {
		return nil, fmt.Errorf("PublishAsset failed. Asset cannot be nil!")
	}

	c.rl.Wait()

	published = new(Asset)
	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets/%v/published", asset.Space.ID, asset.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", asset.System.Version)).
		Receive(published, contentfulError)

	return published, handleError(err, contentfulError)
}

// UnpublishAsset will make the asset unavailable via the Content Delivery API.
func (c *Client) UnpublishAsset(asset *Asset) (unpublished *Asset, err error) {
	if asset == nil {
		return nil, fmt.Errorf("UnpublishAsset failed. Asset cannot be nil!")
	}

	c.rl.Wait()

	unpublished = new(Asset)
	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets/%v/published", asset.Space.ID, asset.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(unpublished, contentfulError)

	return unpublished, handleError(err, contentfulError)
}

// DeleteAsset will delete the specified asset
func (c *Client) DeleteAsset(asset *Asset) (err error) {
	if asset == nil {
		return fmt.Errorf("DeleteAsset failed. Asset cannot be nil!")
	}

	c.rl.Wait()

	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets/%v", asset.Space.ID, asset.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(nil, contentfulError)

	return handleError(err, contentfulError)
}

// ArchiveAsset will archive the asset. An asset can only be archived when they
// are not published. If the asset is published you must first unpublish it.
func (c *Client) ArchiveAsset(asset *Asset) (archived *Asset, err error) {
	if asset == nil {
		err = fmt.Errorf("ArchiveAsset failed. Asset cannot be nil!")
		return
	}

	if err = asset.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	archived = new(Asset)
	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets/%v/archived", asset.Space.ID, asset.System.ID)
	_, err = c.sling.New().
		Put(path).
		Receive(archived, contentfulError)

	return archived, handleError(err, contentfulError)
}

// UnarchiveAsset will unarchive the asset.
func (c *Client) UnarchiveAsset(asset *Asset) (unarchived *Asset, err error) {
	if asset == nil {
		err = fmt.Errorf("UnarchiveAsset failed. Asset cannot be nil!")
		return
	}

	if err = asset.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	unarchived = new(Asset)
	contentfulError := new(Error)
	path := fmt.Sprintf("spaces/%v/assets/%v/archived", asset.Space.ID, asset.System.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(unarchived, contentfulError)

	return unarchived, handleError(err, contentfulError)
}
