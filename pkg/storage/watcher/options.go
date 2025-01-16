package watcher

func WithPrefix[T any]() func(*Watcher[T]) {
	return func(w *Watcher[T]) {
		w.withPrefix = true
	}
}

func WithAfterInitialLoad[T any](callback func()) func(*Watcher[T]) {
	return func(w *Watcher[T]) {
		w.afterInitialLoad = callback
	}
}

func WithBeforeUpdate[T any](callback func(entry T)) func(*Watcher[T]) {
	return func(w *Watcher[T]) {
		w.beforeUpdate = callback
	}
}

func WithKeyFunc[T any](kf func(key string) string) func(*Watcher[T]) {
	return func(w *Watcher[T]) {
		w.keyFunc = kf
	}
}
