package store

type store interface {
	InsertURL(shortcode, url string) error
	RetrieveURL(shortcode string) (string, error)
}
