package management

import (
	"fmt"
)

type EntryFields map[string]map[string]interface{}

type NewEntry struct {
	Fields EntryFields `json:"fields"`
}

// Validate validates the entry
func (c *NewEntry) Validate() error {
	if c.Fields == nil || len(c.Fields) == 0 {
		return fmt.Errorf("NewEntry.Fields cannot be empty!")
	}

	return nil
}

// Entry represent textual content in a space. An entry's data adheres to a
// certain content type.
type Entry struct {
	System `json:"sys"`
	Fields EntryFields `json:"fields"`
}

// Validate will validate the entry. An error is returned if the entry
// is not  valid.
func (c *Entry) Validate() error {
	if c.Space == nil || c.Space.ID == "" {
		return fmt.Errorf("Entry must have a valid Space associated with it!")
	}

	if c.System.ID == "" {
		return fmt.Errorf("Entry.System.ID cannot be empty!")
	}

	if c.Fields == nil || len(c.Fields) == 0 {
		return fmt.Errorf("Entry.Fields cannot be empty!")
	}

	return nil
}

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

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries/%v", spaceID, entryID)
	_, err = c.sling.New().Get(path).Receive(entry, contentfulError)

	return entry, handleError(err, contentfulError)
}

// CreateEntry will create a new entry with an ID specified by the user or
// generated by the system
func (c *Client) CreateEntry(entry *NewEntry, contentType *ContentType) (created *Entry, err error) {
	if entry == nil || contentType == nil {
		err = fmt.Errorf("CreateEntry failed, entry and contentType cannot be nil!")
		return
	}

	if err = entry.Validate(); err != nil {
		return
	}

	if err = contentType.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	created = new(Entry)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries", contentType.Space.ID)
	_, err = c.sling.New().
		Post(path).
		Set("X-Contentful-Content-Type", contentType.ID).
		BodyJSON(entry).
		Receive(created, contentfulError)

	return created, handleError(err, contentfulError)
}

// UpdateEntry will update the specified entry with any changes that you have
// made.
func (c *Client) UpdateEntry(entry *Entry) (updated *Entry, err error) {
	if entry == nil {
		err = fmt.Errorf("CreateEntry failed. Entry must not be nil!")
		return
	}

	if err = entry.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries/%v", entry.Space.ID, entry.System.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", entry.System.Version)).
		BodyJSON(entry).
		Receive(updated, contentfulError)

	return updated, handleError(err, contentfulError)
}

// DeleteEntry will delete the specified entry
func (c *Client) DeleteEntry(entryID string, spaceID string) (err error) {

	c.rl.Wait()

	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries/%v", spaceID, entryID)
	_, err = c.sling.New().
		Delete(path).
		Receive(nil, contentfulError)

	return handleError(err, contentfulError)
}

// PublishEntry makes the entry available via the Content Delivery API
func (c *Client) PublishEntry(entry *Entry) (published *Entry, err error) {
	if entry == nil {
		err = fmt.Errorf("PublishEntry failed. Entry must not be nil!")
		return
	}

	if err = entry.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	published = new(Entry)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries/%v/published", entry.Space.ID, entry.System.ID)
	_, err = c.sling.New().
		Put(path).
		Set("X-Contentful-Version", fmt.Sprintf("%v", entry.System.Version)).
		Receive(published, contentfulError)

	return published, handleError(err, contentfulError)
}

// UnpublishEntry makes the entry unavailable via the Content Delivery API
func (c *Client) UnpublishEntry(entry *Entry) (unpublished *Entry, err error) {
	if entry == nil {
		err = fmt.Errorf("UnpublishEntry failed. Entry must not be nil!")
		return
	}

	if err = entry.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	unpublished = new(Entry)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries/%v/published", entry.Space.ID, entry.System.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(unpublished, contentfulError)

	return unpublished, handleError(err, contentfulError)
}

// ArchiveEntry will archive the specified entry. An entry can only be archived
// when it's not published.
func (c *Client) ArchiveEntry(entry *Entry) (archived *Entry, err error) {
	if entry == nil {
		err = fmt.Errorf("PublishEntry failed. Entry must not be nil!")
		return
	}

	if err = entry.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	archived = new(Entry)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries/%v/archived", entry.Space.ID, entry.System.ID)
	_, err = c.sling.New().
		Put(path).
		Receive(archived, contentfulError)

	return archived, handleError(err, contentfulError)
}

// UnarchiveEntry unarchives the specified entry.
func (c *Client) UnarchiveEntry(entry *Entry) (unarchived *Entry, err error) {
	if entry == nil {
		err = fmt.Errorf("UnpublishEntry failed. Entry must not be nil!")
		return
	}

	if err = entry.Validate(); err != nil {
		return
	}

	c.rl.Wait()

	unarchived = new(Entry)
	contentfulError := new(ContentfulError)
	path := fmt.Sprintf("spaces/%v/entries/%v/archived", entry.Space.ID, entry.System.ID)
	_, err = c.sling.New().
		Delete(path).
		Receive(unarchived, contentfulError)

	return unarchived, handleError(err, contentfulError)
}
