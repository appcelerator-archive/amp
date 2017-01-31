package ns

import (
	"fmt"
	"time"

	"github.com/appcelerator/amp/pkg/mq"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
)

type nsmq struct {
	url       string
	clusterID string
	clientID  string
	client    stan.Conn
}

// New returns an NATS streaming implementation of mq.Interface
func New(url string, clusterID string, clientID string) mq.Interface {
	return &nsmq{url: url, clusterID: clusterID, clientID: clientID}
}

// Connect to NATS streaming
func (ns *nsmq) Connect(timeout time.Duration) error {
	nc, err := nats.Connect(ns.url, nats.Timeout(timeout))
	if err != nil {
		return fmt.Errorf("unable to connect to nats streaming: %v", err)
	}
	ns.client, err = stan.Connect(ns.clusterID, ns.clientID, stan.NatsConn(nc), stan.ConnectWait(timeout))
	if err != nil {
		return fmt.Errorf("unable to connect to nats streaming: %v", err)
	}
	return nil
}

// Close the client
func (ns *nsmq) Close() error {
	return ns.client.Close()
}

func (ns *nsmq) Publish(subject string, data proto.Message) error {
	bytes, err := proto.Marshal(data)
	if err != nil {
		return fmt.Errorf("Cannot marshal data:", err)
	}
	return ns.PublishRaw(subject, bytes)
}

func (ns *nsmq) PublishAsync(subject string, data proto.Message, ah mq.AckHandler) (string, error) {
	bytes, err := proto.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Cannot marshal data:", err)
	}
	return ns.PublishAsyncRaw(subject, bytes, ah)
}

func (ns *nsmq) PublishRaw(subject string, data []byte) error {
	return ns.client.Publish(subject, data)
}

func (ns *nsmq) PublishAsyncRaw(subject string, data []byte, ah mq.AckHandler) (string, error) {
	return ns.client.PublishAsync(subject, data, func(str string, err error) {
		if ah != nil {
			ah(str, err)
		}
	})
}

// A subscription represents a subscription to a stan cluster.
type subscription struct {
	subscription stan.Subscription
}

// Unsubscribe implements the Subscription interface
func (sub *subscription) Unsubscribe() error {
	return sub.subscription.Unsubscribe()
}

// Close implements the Subscription interface
func (sub *subscription) Close() error {
	return sub.subscription.Close()
}

// DefaultSubscriptionOptions are the default subscriptions' options
var DefaultSubscriptionOptions = mq.SubscriptionOptions{
	MaxInflight: stan.DefaultMaxInflight,
	AckWait:     stan.DefaultAckWait,
}

// Convert MQ option closures to NATS option closures
func (ns *nsmq) convertOptions(mqOpts []mq.SubscriptionOption) (stanOpts []stan.SubscriptionOption) {
	options := DefaultSubscriptionOptions
	for _, mqOpt := range mqOpts {
		mqOpt(&options)
	}
	stanOpts = append(stanOpts, stan.DurableName(options.DurableName))
	stanOpts = append(stanOpts, stan.MaxInflight(options.MaxInflight))
	stanOpts = append(stanOpts, stan.AckWait(options.AckWait))
	switch options.StartAt {
	case mq.StartPosition_NewOnly:
	// Nothing to do
	case mq.StartPosition_LastReceived:
		stanOpts = append(stanOpts, stan.StartWithLastReceived())
	case mq.StartPosition_TimeDeltaStart:
		stanOpts = append(stanOpts, stan.StartAtTime(options.StartTime))
	case mq.StartPosition_SequenceStart:
		stanOpts = append(stanOpts, stan.StartAtSequence(options.StartSequence))
	case mq.StartPosition_First:
		stanOpts = append(stanOpts, stan.DeliverAllAvailable())
	}
	if options.ManualAcks {
		stanOpts = append(stanOpts, stan.SetManualAckMode())
	}
	return
}

func (ns *nsmq) Subscribe(subject string, mh mq.MsgHandler, template proto.Message, opts ...mq.SubscriptionOption) (mq.Subscription, error) {
	return ns.SubscribeRaw(subject, func(msg []byte) {
		if mh != nil {
			// create a new empty message from the typed instance template
			val := proto.Clone(template)
			// unmarshall the bytes into the new instance
			err := proto.Unmarshal(msg, val)
			if err != nil {
				mh(nil, err)
			} else {
				mh(val, nil)
			}
		}
	}, opts...)
}

func (ns *nsmq) QueueSubscribe(subject, qgroup string, mh mq.MsgHandler, template proto.Message, opts ...mq.SubscriptionOption) (mq.Subscription, error) {
	return ns.QueueSubscribeRaw(subject, qgroup, func(msg []byte) {
		if mh != nil {
			// create a new empty message from the typed instance template
			val := proto.Clone(template)
			// unmarshall the bytes into the new instance
			err := proto.Unmarshal(msg, val)
			if err != nil {
				mh(nil, err)
			} else {
				mh(val, nil)
			}
		}
	}, opts...)
}

func (ns *nsmq) SubscribeRaw(subject string, mh mq.MsgHandlerRaw, opts ...mq.SubscriptionOption) (mq.Subscription, error) {
	var err error
	sub := &subscription{}
	sub.subscription, err = ns.client.Subscribe(subject, func(msg *stan.Msg) {
		if mh != nil {
			mh(msg.Data)
		}
	}, ns.convertOptions(opts)...)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (ns *nsmq) QueueSubscribeRaw(subject, qgroup string, mh mq.MsgHandlerRaw, opts ...mq.SubscriptionOption) (mq.Subscription, error) {
	var err error
	sub := &subscription{}
	sub.subscription, err = ns.client.QueueSubscribe(subject, qgroup, func(msg *stan.Msg) {
		if mh != nil {
			mh(msg.Data)
		}
	}, ns.convertOptions(opts)...)
	if err != nil {
		return nil, err
	}
	return sub, nil
}
