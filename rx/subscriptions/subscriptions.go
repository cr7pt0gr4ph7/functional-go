package subscriptions

// Same as `rx.Subsctiption`.
type Subscription interface {
	Cancel()
}

type NopSubscription struct{}

var nopInstance Subscription = NopSubscription{}

func Nop() Subscription {
	return nopInstance
}

func (_ NopSubscription) Cancel() {
	// Do nothing, as the name NopSubscription implies.
}
