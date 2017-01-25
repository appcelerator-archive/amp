package mq

import (
	"github.com/golang/protobuf/proto"
	"time"
)

// Interface must be implemented for a message queuer
type Interface interface {
	// Connect to message queuer
	Connect(timeout time.Duration) error

	// Close connection
	Close() error

	// Publish
	Publish(subject string, data proto.Message) error

	// PublishAsync
	PublishAsync(subject string, data proto.Message, ah AckHandler) (string, error)

	// Subscribe
	Subscribe(subject string, mh MsgHandler, template proto.Message, opts ...SubscriptionOption) (Subscription, error)

	// QueueSubscribe
	QueueSubscribe(subject, qgroup string, mh MsgHandler, template proto.Message, opts ...SubscriptionOption) (Subscription, error)

	// PublishRaw
	PublishRaw(subject string, data []byte) error

	// PublishAsyncRaw
	PublishAsyncRaw(subject string, data []byte, ah AckHandler) (string, error)

	// SubscribeRaw
	SubscribeRaw(subject string, mh MsgHandlerRaw, opts ...SubscriptionOption) (Subscription, error)

	// QueueSubscribeRaw
	QueueSubscribeRaw(subject, qgroup string, mh MsgHandlerRaw, opts ...SubscriptionOption) (Subscription, error)
}

// AckHandler is used for Async Publishing to provide status of the ack.
// The func will be passed the GUID and any error state. No error means the
// message was successfully received by message queuer.
type AckHandler func(string, error)

// MsgHandler is a callback function that processes messages delivered to
// asynchronous subscribers.
type MsgHandler func(proto.Message, error)

// MsgHandlerRaw is a callback function that processes messages delivered to
// asynchronous subscribers.
type MsgHandlerRaw func([]byte)

// Enum for start position type.
type StartPosition int32

const (
	StartPosition_NewOnly        StartPosition = 0
	StartPosition_LastReceived   StartPosition = 1
	StartPosition_TimeDeltaStart StartPosition = 2
	StartPosition_SequenceStart  StartPosition = 3
	StartPosition_First          StartPosition = 4
)

// SubscriptionOption is a function on the options for a subscription.
type SubscriptionOption func(*SubscriptionOptions) error

// SubscriptionOptions are used to control the Subscription's behavior.
type SubscriptionOptions struct {
	// DurableName, if set will survive client restarts.
	DurableName string
	// Controls the number of messages the cluster will have inflight without an ACK.
	MaxInflight int
	// Controls the time the cluster will wait for an ACK for a given message.
	AckWait time.Duration
	// StartPosition enum from proto.
	StartAt StartPosition
	// Optional start sequence number.
	StartSequence uint64
	// Optional start time.
	StartTime time.Time
	// Option to do Manual Acks
	ManualAcks bool
}

// Subscription represents a subscription within the message queuer cluster. Subscriptions
// will be rate matched and follow at-least delivery semantics.
type Subscription interface {
	// Unsubscribe removes interest in the subscription.
	// For durables, it means that the durable interest is also removed from
	// the server. Restarting a durable with the same name will not resume
	// the subscription, it will be considered a new one.
	Unsubscribe() error

	// Close removes this subscriber from the server, but unlike Unsubscribe(),
	// the durable interest is not removed. If the client has connected to a server
	// for which this feature is not available, Close() will return an error.
	Close() error
}

// MaxInflight is an Option to set the maximum number of messages the cluster will send
// without an ACK.
func MaxInflight(m int) SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.MaxInflight = m
		return nil
	}
}

// AckWait is an Option to set the timeout for waiting for an ACK from the cluster's
// point of view for delivered messages.
func AckWait(t time.Duration) SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.AckWait = t
		return nil
	}
}

// StartAt sets the desired start position for the message stream.
func StartAt(sp StartPosition) SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.StartAt = sp
		return nil
	}
}

// StartAtSequence sets the desired start sequence position and state.
func StartAtSequence(seq uint64) SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.StartAt = StartPosition_SequenceStart
		o.StartSequence = seq
		return nil
	}
}

// StartAtTime sets the desired start time position and state.
func StartAtTime(start time.Time) SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.StartAt = StartPosition_TimeDeltaStart
		o.StartTime = start
		return nil
	}
}

// StartAtTimeDelta sets the desired start time position and state using the delta.
func StartAtTimeDelta(ago time.Duration) SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.StartAt = StartPosition_TimeDeltaStart
		o.StartTime = time.Now().Add(-ago)
		return nil
	}
}

// StartWithLastReceived is a helper function to set start position to last received.
func StartWithLastReceived() SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.StartAt = StartPosition_LastReceived
		return nil
	}
}

// DeliverAllAvailable will deliver all messages available.
func DeliverAllAvailable() SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.StartAt = StartPosition_First
		return nil
	}
}

// SetManualAckMode will allow clients to control their own acks to delivered messages.
func SetManualAckMode() SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.ManualAcks = true
		return nil
	}
}

// DurableName sets the DurableName for the subscriber.
func DurableName(name string) SubscriptionOption {
	return func(o *SubscriptionOptions) error {
		o.DurableName = name
		return nil
	}
}
