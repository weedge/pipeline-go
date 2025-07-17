package notifiers

// Notifier defines an interface for components that can notify and be waited on.
type Notifier interface {
	Notify()
	Wait() <-chan struct{}
}

// ChannelNotifier is a concrete implementation of Notifier using a Go channel.
type ChannelNotifier struct {
	ch chan struct{}
}

// NewChannelNotifier creates a new ChannelNotifier.
func NewChannelNotifier() *ChannelNotifier {
	return &ChannelNotifier{
		// Buffered channel of size 1, so that a notification can be sent even if no one is waiting.
		ch: make(chan struct{}, 1),
	}
}

// Notify sends a notification. It won't block if the channel is full.
func (n *ChannelNotifier) Notify() {
	select {
	case n.ch <- struct{}{}:
	default: // If the channel is already full, do nothing.
	}
}

// Wait returns a channel that can be waited on.
// A value will be received when a notification is sent.
func (n *ChannelNotifier) Wait() <-chan struct{} {
	return n.ch
}
