package nats

import (
    "context"
    jsoniter "github.com/json-iterator/go"
    "github.com/nats-io/nats.go"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/messaging"
    "regexp"
    "strconv"
    "strings"
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

    config       *messaging.Config
    streamName   string
    subjectNames []string
    pullNumber   int32
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
            n.streamName = config.Nats.JetStream
            n.subjectNames = config.Nats.TopicNames
            if len(n.subjectNames) == 0 {
                n.subjectNames = []string{n.streamName + ".*"}
            }
            n.createStream(n.streamName, n.subjectNames)

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
        if _, err = n.js.AddStream(&nats.StreamConfig{Name: name, Subjects: subjects}); err != nil {
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
                if _, err := n.js.UpdateStream(&nats.StreamConfig{Name: n.streamName, Subjects: allSubjects}); err != nil {
                    logs.Warnw("failed to update stream", "name", n.streamName, "subjects", allSubjects, "error", err.Error())
                }
            }
        }
    }
}

func (n *Nats) updateSubjectName(name string) {
    if len(n.streamName) > 0 && len(name) > 0 {
        for _, s := range n.subjectNames {
            if s == name || s == n.streamName+".*" {
                return
            }
        }

        n.subjectNames = append(n.subjectNames, name)
        if _, err := n.js.UpdateStream(&nats.StreamConfig{Name: n.streamName, Subjects: n.subjectNames}); err != nil {
            logs.Warnw("failed to update stream", "name", n.streamName, "subject", name, "error", err.Error())
            n.subjectNames = n.subjectNames[0 : len(n.subjectNames)-1]
        }
    }
}

func (n *Nats) GetConn() *nats.Conn {
    if n != nil {
        return n.conn
    }
    return nil
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
            m.Ack()
        })
        msg.SetNak(func() {
            m.Nak()
        })
        msg.SetTerm(func() {
            m.Term()
        })
        msg.SetInProgress(func() {
            m.InProgress()
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

func (n *Nats) queueSubscribe(opts *messaging.Subscription, cb nats.MsgHandler) (*nats.Subscription, error) {
    if n.js != nil {
        n.updateSubjectName(opts.Topic)
        if opts.Pull {
            return n.pullSubscribe(opts, cb)
        } else {
            return n.js.QueueSubscribe(opts.Topic, opts.Group, cb)
        }
    } else {
        return n.conn.QueueSubscribe(opts.Topic, opts.Group, cb)
    }
}

func (n *Nats) subscribe(opts *messaging.Subscription, cb nats.MsgHandler) (*nats.Subscription, error) {
    if n.js != nil {
        n.updateSubjectName(opts.Topic)
        return n.js.Subscribe(opts.Topic, cb)
    } else {
        return n.conn.Subscribe(opts.Topic, cb)
    }
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
    //ss := strconv.Itoa(rand.Intn(10000))
    durableName = strings.Join([]string{durableName, suffix}, "-")

    subOpts := []nats.SubOpt{nats.PullMaxWaiting(128)}
    if opts.AckTimeout != nil {
        subOpts = append(subOpts, nats.AckWait(opts.AckTimeout.ToDuration()))
    }
    subs, err := n.js.PullSubscribe(opts.Topic, durableName, subOpts...)
    if err != nil {
        logs.Warnw("failed to pull subscribe the topic", "topic", opts.Topic, "error", err.Error())
        return nil, err
    }

    subs.SetPendingLimits(1000, 5*1024*1024)

    go func(opts *messaging.Subscription, durableName string) {
        logs.Infof("abfuzz.AgentServer subscribed the %s topic", opts.Topic)
        for {
            select {
            case <-n.shutdownCh:
                logs.Infow("shutdown the nats client, will closed the pull subscribe", "subscribe", durableName)
                return
            default:
            }
            logs.Infow("fetching the next message", "subscribe", durableName)
            msgs, _ := subs.Fetch(1)
            for _, msg := range msgs {
                logs.Infow("fetched the message", "subscribe", durableName, "subject", msg.Subject)
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
