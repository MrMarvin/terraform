package state

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"strings"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/states/statemgr"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/version"
)

var rngSource *rand.Rand

func init() {
	rngSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// State is the collection of all state interfaces.
type State interface {
	StateReader
	StateWriter
	StateRefresher
	StatePersister
	Locker
}

// StateReader is the interface for things that can return a state. Retrieving
// the state here must not error. Loading the state fresh (an operation that
// can likely error) should be implemented by RefreshState. If a state hasn't
// been loaded yet, it is okay for State to return nil.
//
// Each caller of this function must get a distinct copy of the state, and
// it must also be distinct from any instance cached inside the reader, to
// ensure that mutations of the returned state will not affect the values
// returned to other callers.
type StateReader interface {
	State() *terraform.State
}

// StateWriter is the interface that must be implemented by something that
// can write a state. Writing the state can be cached or in-memory, as
// full persistence should be implemented by StatePersister.
//
// Implementors that cache the state in memory _must_ take a copy of it
// before returning, since the caller may continue to modify it once
// control returns. The caller must ensure that the state instance is not
// concurrently modified _during_ the call, or behavior is undefined.
//
// If an object implements StatePersister in conjunction with StateReader
// then these methods must coordinate such that a subsequent read returns
// a copy of the most recent write, even if it has not yet been persisted.
type StateWriter interface {
	WriteState(*terraform.State) error
}

// StateRefresher is the interface that is implemented by something that
// can load a state. This might be refreshing it from a remote location or
// it might simply be reloading it from disk.
type StateRefresher interface {
	RefreshState() error
}

// StatePersister is implemented to truly persist a state. Whereas StateWriter
// is allowed to perhaps be caching in memory, PersistState must write the
// state to some durable storage.
//
// If an object implements StatePersister in conjunction with StateReader
// and/or StateRefresher then these methods must coordinate such that
// subsequent reads after a persist return an updated value.
type StatePersister interface {
	PersistState() error
}

// Locker is implemented to lock state during command execution.
// The info parameter can be recorded with the lock, but the
// implementation should not depend in its value. The string returned by Lock
// is an ID corresponding to the lock acquired, and must be passed to Unlock to
// ensure that the correct lock is being released.
//
// Lock and Unlock may return an error value of type LockError which in turn
// can contain the LockInfo of a conflicting lock.
type Locker interface {
	Lock(info *LockInfo) (string, error)
	Unlock(id string) error
}

// test hook to verify that LockWithContext has attempted a lock
var postLockHook func()

// Lock the state, using the provided context for timeout and cancellation.
// This backs off slightly to an upper limit.
func LockWithContext(ctx context.Context, s State, info *LockInfo) (string, error) {
	delay := time.Second
	maxDelay := 16 * time.Second
	for {
		id, err := s.Lock(info)
		if err == nil {
			return id, nil
		}

		le, ok := err.(*LockError)
		if !ok {
			// not a lock error, so we can't retry
			return "", err
		}

		if le == nil || le.Info == nil || le.Info.ID == "" {
			// If we dont' have a complete LockError, there's something wrong with the lock
			return "", err
		}

		if postLockHook != nil {
			postLockHook()
		}

		// there's an existing lock, wait and try again
		select {
		case <-ctx.Done():
			// return the last lock error with the info
			return "", err
		case <-time.After(delay):
			if delay < maxDelay {
				delay *= 2
			}
		}
	}
}

// Generate a LockInfo structure, populating the required fields.
func NewLockInfo() *LockInfo {
	// this doesn't need to be cryptographically secure, just unique.
	// Using math/rand alleviates the need to check handle the read error.
	// Use a uuid format to match other IDs used throughout Terraform.
	buf := make([]byte, 16)
	rngSource.Read(buf)

	id, err := uuid.FormatUUID(buf)
	if err != nil {
		// this of course shouldn't happen
		panic(err)
	}

	// don't error out on user and hostname, as we don't require them
	userName := ""
	if userInfo, err := user.Current(); err == nil {
		userName = userInfo.Username
	}
	host, _ := os.Hostname()

	info := &LockInfo{
		ID:      id,
		Who:     fmt.Sprintf("%s@%s", userName, host),
		Version: version.Version,
		Created: time.Now().UTC(),
	}
	return info
}

type LockInfo = statemgr.LockInfo

type LockError struct {
	Info *LockInfo
	Err  error
}

func (e *LockError) Error() string {
	var out []string
	if e.Err != nil {
		out = append(out, e.Err.Error())
	}

	if e.Info != nil {
		out = append(out, e.Info.String())
	}
	return strings.Join(out, "\n")
}
