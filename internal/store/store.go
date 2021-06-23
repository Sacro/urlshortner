package store

// store is a generic interface for inserting/retrieving shortcodes/urls
type store interface {
	// InsertURL puts the shortcode and URL into the store
	InsertURL(shortcode, url string) error

	// RetrieveURL returns the URL for a shortcode, or an error if not found
	RetrieveURL(shortcode string) (string, error)
}
