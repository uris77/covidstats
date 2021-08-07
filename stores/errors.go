package stores

import "fmt"

type MongoConnectionErr struct {
	Reason string
	Inner  error
}

func (e MongoConnectionErr) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("failed to connect to mongo: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("failed to connect to mongo: %s", e.Reason)
}

func (e MongoConnectionErr) Unwrap() error {
	return e.Inner
}

// Errors when doing queries
type MongoQueryErr struct {
	Reason string
	Inner  error
}

func (e MongoQueryErr) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("failed to query mongo: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("failed to query mongo: %s", e.Reason)
}

type MongoNoResultErr struct {
	Reason string
	Inner  error
}

func (e MongoNoResultErr) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("mongo: no records returned from query: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("mongo: no records returned from query: %s", e.Reason)
}

func (e MongoNoResultErr) Unwrap() error {
	return e.Inner
}
