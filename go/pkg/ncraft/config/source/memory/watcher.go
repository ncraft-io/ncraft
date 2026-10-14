package memory

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
)

type watcher struct {
	Id      string
	Updates chan *source.ChangeSet
	Source  *memory
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	cs, ok := <-w.Updates
	if !ok {
		return nil, source.ErrWatcherStopped
	}
	return cs, nil
}

func (w *watcher) Stop() error {
	w.Source.Lock()
	if _, ok := w.Source.Watchers[w.Id]; ok {
		delete(w.Source.Watchers, w.Id)
		close(w.Updates)
	}
	w.Source.Unlock()
	return nil
}
