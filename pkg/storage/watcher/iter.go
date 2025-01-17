package watcher

import (
	"iter"
	"strings"
)

func (w *Watcher[T]) iter(fullkey bool) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		w.mutex.RLock()
		defer w.mutex.RUnlock()
		for k, v := range w.entries {
			kk := k
			if !fullkey {
				kk = strings.TrimPrefix(kk, w.prefix.String())
			}
			if !yield(kk, v) {
				return
			}
		}
	}
}

func (w *Watcher[T]) Iter() iter.Seq2[string, T] {
	return w.iter(true)
}

func (w *Watcher[T]) IterRelativeKey() iter.Seq2[string, T] {
	return w.iter(false)
}
