package util

// Validator lets all objects that implement the Validate() method to be accepted as the same type
type Validator interface {
	Validate() error
}
