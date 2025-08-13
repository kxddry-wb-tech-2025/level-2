package models

import "net/url"

// ResourceType is used to differentiate types
type ResourceType int

const (
	// HTML means "this file is HTML"
	HTML ResourceType = iota
	// CSS means the same
	CSS
	// JS means this is a script written in JavaScript
	JS
	// Image means this is an image
	Image
	// Other is used for unsupported other types.
	Other
)

// Resource is used to mark out a resource that needs to be downloaded
type Resource struct {
	URL  *url.URL
	Type ResourceType
}

// Task is a helper structure for resources
type Task struct {
	URL   *url.URL
	Depth int
	Type  ResourceType
}
