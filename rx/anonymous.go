package rx

type AnonymousObservable[T any] func(observer Observer[T]) Subscription

func NewObservable[T any](onSubscribe AnonymousObservable[T]) Observable[T] {
	return onSubscribe
}

func (o AnonymousObservable[T]) Subscribe(observer Observer[T]) Subscription {
	return o(observer)
}

type AnonymousObserverConfig[T any] struct {
	Next  func(value T)
	Done  func()
	Error func(err error)
}

type AnonymousObserver[T any] struct {
	config AnonymousObserverConfig[T]
}

func NewObserver[T any](funcs AnonymousObserverConfig[T]) Observer[T] {
	return AnonymousObserver[T]{config: funcs}
}

func (o AnonymousObserver[T]) Next(value T)    { o.config.Next(value) }
func (o AnonymousObserver[T]) Done()           { o.config.Done() }
func (o AnonymousObserver[T]) Error(err error) { o.config.Error(err) }

type AnonymousSubscription func()

func NewSubscription(onCancel AnonymousSubscription) Subscription {
	return onCancel
}

func (s AnonymousSubscription) Cancel() {
	s()
}
