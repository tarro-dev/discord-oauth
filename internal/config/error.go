package config

// ErrMissingKey is an error returned when a key is missing.
type ErrMissingKey struct {
	Key string
}

// Error returns the error message.
func (e ErrMissingKey) Error() string {
	return "missing key: " + e.Key
}
