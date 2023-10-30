package nats

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/messaging"
)

const (
	defaultPullMaxWaiting = 128
)

func init() {
	messaging.Register("nats", func(config *messaging.Config) (messaging.Queue, error) {
		return New(config)
	})
}

type Msg = nats.Msg
type Subscription = nats.Subscription

var ErrConnectionClosed = nats.ErrConnectionClosed
var ErrTimeout = nats.ErrTimeout
var ErrNoResponders = nats.ErrNoResponders

var PullMaxWaiting = nats.PullMaxWaiting
var Context = nats.Context

// Nats implement the Queue interface
type Nats struct {
	id            string
	conn          *nats.Conn
	js            nats.JetStreamContext
	subscriptions map[string]*nats.Subscription
	shutdown      bool
	shutdownCh    chan struct{}

	config     *messaging.Config
	pullNumber int32
}

// New NewNats creates a new Nats connection
func New(config *messaging.Config) (*Nats, error) {
	if nc, err := nats.Connect(config.ServiceName); err != nil {
		return nil, err
	} else {
		n := &Nats{
			conn:          nc,
			id:            "",
			subscriptions: map[string]*nats.Subscription{},
		}

		if config.Nats != nil && len(config.Nats.JetStream) > 0 {
			js, err := nc.JetStream()
			if err != nil {
				return nil, err
			}
			n.js = js
			if len(n.SubjectNames()) == 0 {
				n.SetSubjectNames(n.StreamName() + ".*")
			}
			n.createStream(n.StreamName(), n.SubjectNames())
			n.shutdownCh = make(chan struct{})
		}

		return n, nil
	}
}

func (n *Nats) createStream(name string, subjects []string) {
	stream, err := n.js.StreamInfo(name)
	if err != nil {
		logs.Warnw("failed to get the stream info", "error", err.Error())
	}

	if stream == nil {
		streamConfig := &nats.StreamConfig{Name: name, Subjects: subjects}
		if n.MaxMsgs() > 0 {
			streamConfig.MaxMsgs = n.MaxMsgs()
		}
		if n.MaxAge() > 0 {
			streamConfig.MaxAge = time.Duration(n.MaxAge()) * time.Second
		}
		if _, err = n.js.AddStream(streamConfig); err != nil {
			logs.Warnw("failed to creating stream", "name", name, "subjects", subjects, "error", err.Error())
		} else {
			logs.Infof("creating stream %q and subject %q", name, subjects)
		}
	} else {
		if len(subjects) > 0 {
			allSubjects := stream.Config.Subjects
			if len(allSubjects) == 0 {
				allSubjects = subjects
			} else {
				var nonExists []string
				for _, s := range subjects {
					exist := false
					for _, sub := range allSubjects {
						if sub == s {
							exist = true
							break
						}
					}
					if !exist {
						nonExists = append(nonExists, s)
					}
				}
				allSubjects = append(allSubjects, nonExists...)

				streamConfig := &nats.StreamConfig{Name: name, Subjects: allSubjects}
				if n.MaxMsgs() > 0 {
					streamConfig.MaxMsgs = n.MaxMsgs()
				}
				if n.MaxAge() > 0 {
					streamConfig.MaxAge = time.Duration(n.MaxAge()) * time.Second
				}
				if _, err = n.js.UpdateStream(streamConfig); err != nil {
					logs.Warnw("failed to update stream", "name", n.StreamName(), "subjects", allSubjects, "error", err.Error())
				}
			}
		}
	}
}

func (n *Nats) updateSubjectName(name string) {
	if len(n.StreamName()) > 0 && len(name) > 0 {
		for _, s := range n.SubjectNames() {
			if s == name || s == n.StreamName()+".*" {
				return
			}
		}

		n.AppendSubjectName(name)
		if _, err := n.js.UpdateStream(&nats.StreamConfig{Name: n.StreamName(), Subjects: n.SubjectNames()}); err != nil {
			logs.Warnw("failed to update stream", "name", n.StreamName(), "subject", name, "error", err.Error())
			n.SetSubjectNames(n.SubjectNames()[0 : len(n.SubjectNames())-1]...)
		}
	}
}

func (n *Nats) GetConn() *nats.Conn {
	if n != nil {
		return n.conn
	}
	return nil
}

func (n *Nats) StreamName() string {
	if n != nil && n.config != nil {
		return n.config.Nats.GetJetStream()
	}
	return ""
}

func (n *Nats) SubjectNames() []string {
	if n != nil && n.config != nil {
		return n.config.Nats.GetTopicNames()
	}
	return nil
}

func (n *Nats) AppendSubjectName(name string) {
	if n != nil && n.config != nil {
		n.config.Nats.TopicNames = append(n.config.Nats.TopicNames, name)
	}
}

func (n *Nats) SetSubjectNames(names ...string) {
	if n != nil && n.config != nil {
		n.config.Nats.TopicNames = names
	}
}

func (n *Nats) MaxMsgs() int64 {
	if n != nil && n.config != nil {
		return n.config.Nats.GetMaxMsgs()
	}
	return 0
}

func (n *Nats) MaxAge() int64 {
	if n != nil && n.config != nil {
		return n.config.Nats.GetMaxAge()
	}
	return 0
}

func (n *Nats) GetJetStream() nats.JetStreamContext {
	if n != nil {
		return n.js
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

		if err = n.publish(ctx, topic, b); err != nil {
			return logs.NewErrorw("couldn't publish to nats", "error", err)
		}
	}

	logs.Infow("Nats: Published to queue", "topic", topic, "id", n.id)
	return nil
}

func (n *Nats) publish(ctx context.Context, topic string, msg []byte) (err error) {
	_ = ctx
	if n.js != nil {
		n.updateSubjectName(topic)
		_, err = n.js.Publish(topic, msg)
	} else {
		err = n.conn.Publish(topic, msg)
	}
	return
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
			err := m.Ack()
			logs.Warnw("Nats: failed to Ack", "topic", opts.Topic, "error", err)
		})
		msg.SetNak(func() {
			err := m.Nak()
			logs.Warnw("Nats: failed to Nak", "topic", opts.Topic, "error", err)
		})
		msg.SetTerm(func() {
			err := m.Term()
			logs.Warnw("Nats: failed to Term", "topic", opts.Topic, "error", err)
		})
		msg.SetInProgress(func() {
			err := m.InProgress()
			logs.Warnw("Nats: failed to InProgress", "topic", opts.Topic, "error", err)
		})

		err := h(context.Background(), opts, msg)
		if err != nil {
			return
		}

		if opts.AutoAck {
			msg.Ack()
		}
	}

	if len(opts.Group) > 0 {
		sub, err := n.queueSubscribe(opts, subscribe)
		n.subscriptions[opts.Topic] = sub
		if err != nil {
			logs.Warnw("Nats: failed to subscribe with group", "topic", opts.Topic, "group", opts.Group, "error", err)
		}
	} else {
		sub, err := n.subscribe(opts, subscribe)
		n.subscriptions[opts.Topic] = sub
		if err != nil {
			logs.Warnw("Nats: failed to subscribe", "topic", opts.Topic, "error", err)
		}
	}
}

func (n *Nats) queueSubscribe(opts *messaging.Subscription, cb nats.MsgHandler) (sub *nats.Subscription, err error) {
	if n.js != nil {
		n.updateSubjectName(opts.Topic)
		if opts.Pull {
			sub, err = n.pullSubscribe(opts, cb)
		} else {
			sub, err = n.js.QueueSubscribe(opts.Topic, opts.Group, cb)
		}
	} else {
		sub, err = n.conn.QueueSubscribe(opts.Topic, opts.Group, cb)
	}
	if err != nil {
		return nil, err
	}

	err = setPendingLimit(sub, opts)
	return
}

func (n *Nats) subscribe(opts *messaging.Subscription, cb nats.MsgHandler) (sub *nats.Subscription, err error) {
	if n.js != nil {
		n.updateSubjectName(opts.Topic)
		sub, err = n.js.Subscribe(opts.Topic, cb)
	} else {
		sub, err = n.conn.Subscribe(opts.Topic, cb)
	}
	if err != nil {
		return nil, err
	}

	err = setPendingLimit(sub, opts)
	return
}

func setPendingLimit(sub *nats.Subscription, opts *messaging.Subscription) error {
	msgLimit := int(opts.PendingMsgLimit)
	bytesLimit := int(opts.PendingBytesLimit)
	if msgLimit != 0 || bytesLimit != 0 {
		if msgLimit == 0 {
			msgLimit = nats.DefaultSubPendingMsgsLimit
		}
		if bytesLimit == 0 {
			bytesLimit = nats.DefaultSubPendingBytesLimit
		}
		if err := sub.SetPendingLimits(msgLimit, bytesLimit); err != nil {
			logs.Warnw("failed to set pending limits", "name", opts.Name, "topic", opts.Topic,
				"msgLimit", msgLimit, "bytesLimit", bytesLimit, "error", err)
			return err
		}
	}
	return nil
}

var durableNameRegex = regexp.MustCompile("[a-zA-Z0-9_-]+")

func (n *Nats) pullSubscribe(opts *messaging.Subscription, cb nats.MsgHandler) (*nats.Subscription, error) {
	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	durableName := opts.Group
	n.pullNumber++

	suffix := ""
	segments := strings.Split(opts.Topic, ".")
	if len(segments) > 0 {
		suffix = segments[len(segments)-1]
	}
	if !durableNameRegex.MatchString(suffix) {
		suffix = strconv.Itoa(int(n.pullNumber))
	}

	durableName = strings.Join([]string{durableName, suffix}, "-")

	maxWaiting := opts.PullMaxWaiting
	if maxWaiting == 0 {
		maxWaiting = defaultPullMaxWaiting
	}
	subOpts := []nats.SubOpt{nats.PullMaxWaiting(int(maxWaiting))}
	if opts.AckTimeout != nil {
		subOpts = append(subOpts, nats.AckWait(opts.AckTimeout.ToDuration()))
	}
	subs, err := n.js.PullSubscribe(opts.Topic, durableName, subOpts...)
	if err != nil {
		logs.Warnw("failed to pull subscribe the topic", "topic", opts.Topic, "error", err.Error())
		return nil, err
	}

	go func(opts *messaging.Subscription, durableName string) {
		logs.Infow("subscribed ok", "topic", opts.Topic, "subscribe", durableName)
		for {
			select {
			case <-n.shutdownCh:
				logs.Infow("shutdown the nats client, will closed the pull subscribe", "subscribe", durableName)
				return
			default:
			}
			logs.Debugw("fetching the next message", "subscribe", durableName)
			ms, _ := subs.Fetch(1)
			for _, msg := range ms {
				logs.Debugw("fetched the message", "subscribe", durableName, "subject", msg.Subject)
				cb(msg)
			}
		}
	}(opts, durableName)

	return subs, nil
}

// Shutdown shuts down all subscribers
func (n *Nats) Shutdown() {
	n.conn.Close()
	close(n.shutdownCh)
	n.shutdown = true
	return
}
