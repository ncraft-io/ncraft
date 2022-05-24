package nats

import (
    "context"
    jsoniter "github.com/json-iterator/go"
    "github.com/nats-io/nats.go"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/messaging"
)

func init() {
    messaging.Register("nats", func(config *messaging.Config) (messaging.Queue, error) {
        return New(config)
    })
}

type Msg = nats.Msg

var ErrConnectionClosed = nats.ErrConnectionClosed
var ErrTimeout = nats.ErrTimeout
var ErrNoResponders = nats.ErrNoResponders

// Nats implement the Queue interface
type Nats struct {
    id            string
    conn          *nats.Conn
    subscriptions map[string]*nats.Subscription
    shutdown      bool
}

// New NewNats creates a new Nats connection
func New(config *messaging.Config) (*Nats, error) {
    if nc, err := nats.Connect(config.ServiceName); err != nil {
        return nil, err
    } else {
        return &Nats{
            conn:          nc,
            id:            "",
            subscriptions: map[string]*nats.Subscription{},
        }, nil
    }
}

func (n *Nats) GetConn() *nats.Conn {
    if n != nil {
        return n.conn
    }
    return nil
}

// Publish to public the messages to topic
func (n *Nats) Publish(ctx context.Context, topic string, messages ...*messaging.Message) error {
    for _, message := range messages {
        b, err := jsoniter.ConfigFastest.Marshal(message)
        if err != nil {
            return err
        }

        if err = n.conn.Publish(topic, b); err != nil {
            return logs.NewErrorw("couldn't publish to nats", "error", err)
        }
    }

    logs.Infow("Nats: Published to queue", "topic", topic, "id", n.id)
    return nil
}

// Subscribe implements Subscribe for Nats
func (n *Nats) Subscribe(opts *messaging.Subscription, h messaging.Handler) {
    subscribe := func(m *nats.Msg) {
        msg := &messaging.SubMessage{}
        if err := jsoniter.ConfigFastest.Unmarshal(m.Data, msg); err != nil {
            logs.Warnw("Nats: Failed to unmarshal msg from topic", "topic", opts.Topic, "error", err.Error())
            return
        }

        msg.SetAck(func() {
            m.Ack()
        })
        err := h(context.Background(), msg)
        if err != nil {
            return
        }

        if opts.AutoAck {
            msg.Ack()
        }
    }

    if len(opts.Group) > 0 {
        sub, err := n.conn.QueueSubscribe(opts.Topic, opts.Group, subscribe)
        n.subscriptions[opts.Topic] = sub
        if err != nil {
            logs.Warnw("Nats: failed to subscribe with group", "topic", opts.Topic, "group", opts.Group, "error", err)
        }
    } else {
        sub, err := n.conn.Subscribe(opts.Topic, subscribe)
        n.subscriptions[opts.Topic] = sub
        if err != nil {
            logs.Warnw("Nats: failed to subscribe", "topic", opts.Topic, "error", err)
        }
    }
}

// Shutdown shuts down all subscribers
func (n *Nats) Shutdown() {
    n.conn.Close()
    n.shutdown = true
    return
}
