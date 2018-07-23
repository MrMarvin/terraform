package plans

import (
	"sync"
)

// ChangesBuilder is a wrapper around a Changes that provides a concurrency-safe
// interface to insert new changes into the changes.
//
// Each ChangesBuilder is independent of all others, so all concurrent writers
// to a particular Changes must share a single ChangesBuilder. Behavior is
// undefined if any other caller makes changes to the underlying Changes
// object or its nested objects concurrently with any of the methods of a
// particular ChangesBuilder.
type ChangesBuilder struct {
	lock    sync.Mutex
	changes *Changes
}
