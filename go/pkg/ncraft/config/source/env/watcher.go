package env

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
	"sync"
)

type watcher struct {
	exit     chan struct{}
	stopOnce sync.Once
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	<-w.exit

	return nil, source.ErrWatcherStopped
}

func (w *watcher) Stop() error {
	w.stopOnce.Do(func() { close(w.exit) })
	return nil
}

func newWatcher() (source.Watcher, error) {
	return &watcher{exit: make(chan struct{})}, nil
}
