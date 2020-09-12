package gocache

import (
	"time"
)

type work interface {
	do()
}

type job struct {
	cache   *Cache
	key     string
	content interface{}
	life    time.Duration
}

type addingJob struct {
	job
}

type deleteJob struct {
	job
}

func (w *addingJob) do() {
	w.cache.Add(w.key, w.content, w.life)
}

func (w *deleteJob) do() {
	w.cache.DeleteCachedItem(w.key)
}
