package rx

type Observable[T any] interface {
	Subscribe(observer Observer[T]) Subscription
}

type Observer[T any] interface {
	Next(value T)
	Done()
	Error(err error)
}

type Subscription interface {
	Close()
}
