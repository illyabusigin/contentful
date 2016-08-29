package management

import (
	"fmt"
)

// Asset represents a file within a space. An asset can be any kind of file: an
// image, a video, an audio file, a PDF or any other filetype. Assets are
// usually attached to entries through links.
//
// Assets can optionally be localized by providing separate files for each
// locale. Those assets which are not localized simply provide a single file
// under the default locale.
type Asset struct {
	System `json:"sys"`
	Fields AssetFields `json:"fields,omitempty"`
}

// AssetFields contains all asset information.
type AssetFields struct {
	Title  map[string]string    `json:"title,omitempty"`
	File   map[string]AssetData `json:"file,omitempty"`
	Detail *AssetDetail         `json:"details,omitempty"`
}

// AssetData contains all asset information
type AssetData struct {
	MIMEType string `json:"contentType"`
	Name     string `json:"fileName"`
	URL      string `json:"url,omitempty"`
	Upload   string `json:"upload,omitempty"`
}

// Validate will validate the Asset to ensure all necessary fields are present.
func (a *Asset) Validate() error {
	if a.System.ID == "" {
		fmt.Println("system ID is empty!")
		return fmt.Errorf("Asset validation failed. System.ID cannot be empty!")
	}

	if a.System.Space == nil {
		return fmt.Errorf("Asset validation failed. System.Space.ID cannot be empty!")
	}

	return nil
}

// Link represents a link to the Entry and implements the Linkable interface
func (a *Asset) Link() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"sys": map[string]interface{}{
			"id":       a.ID,
			"linkType": "Asset",
			"type":     LinkType,
		},
	}
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

	Fields FileFields `json:"fields"`
}

type FileFields struct {
	Title map[string]string   `json:"title"`
	File  map[string]FileData `json:"file"`
}

// FileData contains all file information
type FileData struct {
	MIMEType string `json:"contentType"`
	Name     string `json:"fileName,omitempty"`
	URL      string `json:"upload,omitempty"`
}

func (f *File) Validate() error {
	if f.SpaceID == "" {
		return fmt.Errorf("Filed validation failed. SpaceID cannot be empty!")
	}

	if f.Fields.File == nil || len(f.Fields.File) == 0 {
		return fmt.Errorf("Filed validation failed. Fields.File cannot be empty!")
	}

	if f.Fields.Title == nil || len(f.Fields.Title) == 0 {
		return fmt.Errorf("Filed validation failed. Fields.Title cannot be empty!")
	}

	for _, data := range f.Fields.File {
		if data.Name == "" {
			return fmt.Errorf("Filed validation failed. FileData.Name cannot be empty. FileData: %v", data)
		} else if data.MIMEType == "" {
			return fmt.Errorf("Filed validation failed. FileData.MIMEType cannot be empty. FileData: %v", data)
		} else if data.URL == "" {
			return fmt.Errorf("Filed validation failed. FileData.URL cannot be empty. FileData: %v", data)
		}
	}

	return nil
}

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
	contentfulError := new(ContentfulError)
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
	contentfulError := new(ContentfulError)
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
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v", spaceID, assetID)
	_, err = c.sling.New().
		Get(path).
		BodyJSON(asset).
		Receive(asset, contentfulError)

	return asset, handleError(err, contentfulError)
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

	contentfulError := new(ContentfulError)
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
	contentfulError := new(ContentfulError)
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
	contentfulError := new(ContentfulError)
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

	contentfulError := new(ContentfulError)
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
	contentfulError := new(ContentfulError)
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
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/assets/%v/archived", asset.Space.ID, asset.System.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(unarchived, contentfulError)

	return unarchived, handleError(err, contentfulError)
}
