package management

import (
	"fmt"
	"time"
)

// Asset represents a file within a space. An asset can be any kind of file: an
// image, a video, an audio file, a PDF or any other filetype. Assets are
// usually attached to entries through links.
//
// Assets can optionally be localized by providing separate files for each
// locale. Those assets which are not localized simply provide a single file
// under the default locale.
type Asset struct {
	System
	Fields struct {
		Title  map[string]string   `json:"title"`
		File   map[string]FileData `json:"file"`
		Detail AssetDetail         `json:"details"`
	} `json:"fields"`

	Processed bool
	Published bool
}

func (a *Asset) Validate() error {
	return nil
}

type ImageDetail struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AssetDetail struct {
	Image *ImageDetail `json:"image,omitempty"`
	Size  int          `json:"size"`
}

// File represents all asset data prior to upload
type File struct {
	SpaceID string `json:"-"`

	Fields struct {
		Title map[string]string   `json:"title"`
		File  map[string]FileData `json:"file"`
	} `json:"fields"`
}

func (f *File) Validate() error {
	return nil
}

// FileData contains all file information
type FileData struct {
	MIMEType string  `json:"contentType"`
	Name     string  `json:"fileName"`
	URL      *string `json:"upload"`
}

// CreateAsset creates a new asset. It's important to note that the asset still
// needs to be processed and published to be availably through the delivery API.
func (c *Client) CreateAsset(file *File) (created *Asset, err error) {
	if file == nil {
		err = fmt.Errorf("CreateAsset failed. Entry must not be nil!")
	}

	if err = file.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets", file.SpaceID)
	_, err = c.sling.New().
		Post(path).
		BodyJSON(file).
		Receive(created, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// UpdateAsset will update the asset
func (c *Client) UpdateAsset(asset *Asset) (updated *Asset, err error) {
	if asset == nil {
		err = fmt.Errorf("UpdateAsset failed. Asset must not be nil!")
	}

	if err = asset.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v", asset.System.Space.ID, asset.System.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", asset.System.Version)).
		BodyJSON(asset).
		Receive(updated, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// FetchAsset will return the specified asset.
func (c *Client) FetchAsset(spaceID string, assetID string) (asset *Asset, err error) {
	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v", spaceID, assetID)
	_, err = c.sling.New().
		Get(path).
		BodyJSON(asset).
		Receive(asset, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// FetchAssets will return all assets associated with a space. You can toggle
// the published flag to only fetch published assets.
func (c *Client) FetchAssets(spaceID string, published bool, limit int, offset int) (assets []*Asset, pagination *Pagination, err error) {
	if spaceID == "" {
		return nil, nil, fmt.Errorf("FetchAssets failed. Space identifier is not valid!")
	}

	if limit <= 0 {
		return nil, nil, fmt.Errorf("FetchAssets failed. Limit must be greater than 0")
	}

	if limit > 100 {
		limit = 100
	}

	c.rl.Wait()

	type assetsResponse struct {
		*Pagination
		Items []*Asset `json:"items"`
	}

	results := new(assetsResponse)
	contentfulError := new(ContentfulError)
	path := func() string {
		if published {
			return fmt.Sprintf("spaces/%v/public/assets", spaceID)
		}

		return fmt.Sprintf("spaces/%v/assets", spaceID)
	}
	req, err := c.sling.New().
		Get(path()).
		Receive(results, contentfulError)

	if err != nil {
		return
	}

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return results.Items, results.Pagination, nil
}

// ProcessAsset process the asset. This uploads the asset to
// Contentful among other things. Processing happens asynchronously, the call
// will not block until it has finished.
func (c *Client) ProcessAsset(asset *Asset, localeCode string) (err error) {
	if asset == nil {
		return fmt.Errorf("ProcessAsset failed. Asset cannot be nil!")
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v/%v/process", asset.Space.ID, asset.ID, localeCode)
	_, err = c.sling.New().
		Put(path).
		Receive(nil, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// PublishAsset will make the asset available via the Content Delivery API.
func (c *Client) PublishAsset(asset *Asset) (published *Asset, err error) {
	if asset == nil {
		return nil, fmt.Errorf("PublishAsset failed. Asset cannot be nil!")
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v/published", asset.Space.ID, asset.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", asset.System.Version)).
		Receive(published, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// UnpublishAsset will make the asset unavailable via the Content Delivery API.
func (c *Client) UnpublishAsset(asset *Asset) (unpublished *Asset, err error) {
	if asset == nil {
		return nil, fmt.Errorf("UnpublishAsset failed. Asset cannot be nil!")
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v/published", asset.Space.ID, asset.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(unpublished, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// DeleteAsset will delete the specified asset
func (c *Client) DeleteAsset(asset *Asset) (err error) {
	if asset == nil {
		return fmt.Errorf("DeleteAsset failed. Asset cannot be nil!")
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v", asset.Space.ID, asset.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(nil, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// ArchiveAsset will archive the asset. An asset can only be archived when they
// are not published. If the asset is published you must first unpublish it.
func (c *Client) ArchiveAsset(asset *Asset) (archived *Asset, err error) {
	if asset == nil {
		err = fmt.Errorf("ArchiveAsset failed. Entry must not be nil!")
	}

	if err = asset.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v/archived", asset.Space.ID, asset.System.ID)
	_, err = c.sling.New().
		Put(path).
		Receive(archived, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}

// UnarchiveAsset will unarchive the asset.
func (c *Client) UnarchiveAsset(asset *Asset) (unarchived *Asset, err error) {
	if asset == nil {
		err = fmt.Errorf("UnarchiveAsset failed. Entry must not be nil!")
	}

	if err = asset.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v/archived", asset.Space.ID, asset.System.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(unarchived, contentfulError)

	if contentfulError.Message != "" {
		err = contentfulError
		return
	}

	return
}
