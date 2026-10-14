package source

import (
	"sync"
)

type noopWatcher struct {
	exit     chan struct{}
	stopOnce sync.Once
}

func (w *noopWatcher) Next() (*ChangeSet, error) {
	<-w.exit

	return nil, ErrWatcherStopped
}

func (w *noopWatcher) Stop() error {
	w.stopOnce.Do(func() { close(w.exit) })
	return nil
}

// NewNoopWatcher returns a watcher that blocks on Next() until Stop() is called.
func NewNoopWatcher() (Watcher, error) {
	return &noopWatcher{exit: make(chan struct{})}, nil
}
