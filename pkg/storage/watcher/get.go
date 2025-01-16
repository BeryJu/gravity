package watcher

func (w *Watcher[T]) Get(key string) T {
	entry, _ := w.GetOK(key)
	return entry
}

func (w *Watcher[T]) GetPrefix(parts ...string) (T, bool) {
	return w.GetOK(w.Prefix().Add(parts...).String())
}

func (w *Watcher[T]) GetOK(key string) (T, bool) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	entry, ok := w.entries[key]
	return entry, ok
}
