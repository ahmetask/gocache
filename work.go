package gocache

import (
	"time"
)

type work interface {
	do()
}

type addingJob struct {
	cache   *Cache
	key     string
	content interface{}
	life    time.Duration
}

func (w *addingJob) do() {
	w.cache.mutex.Lock()

	_, found := w.cache.get(w.key)

	if found {
		w.cache.mutex.Unlock()
		return
	}

	w.cache.set(w.key, w.content, w.life)

	w.cache.mutex.Unlock()
}
