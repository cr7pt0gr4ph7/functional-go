package subscriptions

// Same as `rx.Subsctiption`.
type Subscription interface {
	Cancel()
}

type CancelFunc func()

type Anonymous CancelFunc

func New(onCancel CancelFunc) Subscription {
	return Anonymous(onCancel)
}

func (s Anonymous) Cancel() {
	// TODO(lw) Handle double cancellation by invoking s only once
	if s != nil {
		s()
	}
}

type NopSubscription struct{}

var nopInstance Subscription = NopSubscription{}

func Nop() Subscription {
	return nopInstance
}

func (_ NopSubscription) Cancel() {
	// Do nothing, as the name Nop implies.
}

type State struct {
	cancelled bool
}

func (s *State) Cancel() {
	s.cancelled = true
}

func (s *State) IsCancelled() bool {
	return s.cancelled
}
