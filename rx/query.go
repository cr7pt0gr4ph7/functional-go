package rx

type Query[T any] struct {
	source Observable[T]
}

func NewQuery[T any](source Observable[T]) Query[T] {
	return Query[T]{source: source}
}

func (q Query[T]) Subscribe(observer Observer[T]) Subscription {
	return q.source.Subscribe(observer)
}
