package stores

import "fmt"

// MongoConnectionErr is the error generated when connecting to Mongo fails
type MongoConnectionErr struct {
	Reason string
	Inner  error
}

// Error satisfies the error interface
func (e MongoConnectionErr) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("failed to connect to mongo: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("failed to connect to mongo: %s", e.Reason)
}

// Unwrap satisfies the error interface
func (e MongoConnectionErr) Unwrap() error {
	return e.Inner
}

// MongoQueryErr are the errors when doing queries
type MongoQueryErr struct {
	Reason string
	Inner  error
}

// Error satisfies the error interface
func (e MongoQueryErr) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("failed to query mongo: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("failed to query mongo: %s", e.Reason)
}

// MongoNoResultErr is the error when no results are returned from a mongo query
type MongoNoResultErr struct {
	Reason string
	Inner  error
}

// Error satisfies the error interface
func (e MongoNoResultErr) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("mongo: no records returned from query: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("mongo: no records returned from query: %s", e.Reason)
}

// Unwrap satisfies the error interface
func (e MongoNoResultErr) Unwrap() error {
	return e.Inner
}
