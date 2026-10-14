//go:build !linux
// +build !linux

package file

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	f *file

	fw       *fsnotify.Watcher
	exit     chan bool
	stopOnce sync.Once
	stopErr  error
}

func newWatcher(f *file) (source.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := fw.Add(f.path); err != nil {
		_ = fw.Close()
		return nil, err
	}

	return &watcher{
		f:    f,
		fw:   fw,
		exit: make(chan bool),
	}, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	// is it closed?
	select {
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	default:
	}

	// try get the event
	select {
	case event, ok := <-w.fw.Events:
		if !ok {
			return nil, source.ErrWatcherStopped
		}
		if event.Op == fsnotify.Rename {
			// check existence of file, and add watch again
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		}

		c, err := w.f.Read()
		if err != nil {
			return nil, err
		}
		return c, nil
	case err, ok := <-w.fw.Errors:
		if !ok {
			return nil, source.ErrWatcherStopped
		}
		return nil, err
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	}
}

func (w *watcher) Stop() error {
	w.stopOnce.Do(func() {
		close(w.exit)
		w.stopErr = w.fw.Close()
	})
	return w.stopErr
}
