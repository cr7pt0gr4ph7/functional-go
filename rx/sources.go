package rx

import (
	"github.com/cr7pt0gr4ph7/functional-go/rx/subscriptions"
)

func Never[T any]() Observable[T] {
	return NewObservable(func(_ Observer[T]) Subscription {
		// The observer is never notified
		return subscriptions.Nop()
	})
}

func Fail[T any](err error) Observable[T] {
	return NewObservable(func(observer Observer[T]) Subscription {
		observer.Error(err)
		return subscriptions.Nop()
	})
}

func Completed[T any]() Observable[T] {
	return NewObservable(func(observer Observer[T]) Subscription {
		observer.Done()
		return subscriptions.Nop()
	})
}

func Forever[T any](value T) Observable[T] {
	return NewObservable(func(observer Observer[T]) Subscription {
		// TODO(lw) It is currently impossible to actually cancel the subscription
		var s subscriptions.State
		for !s.IsCancelled() {
			observer.Next(value)
		}
		return s
	})
}

func Return[T any](value T) Observable[T] {
	return NewObservable(func(observer Observer[T]) Subscription {
		observer.Next(value)
		observer.Done()
		return subscriptions.Nop()
	})
}
