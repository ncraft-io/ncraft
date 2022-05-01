package messaging

import "github.com/mojo-lang/core/go/pkg/mojo/core"

type SubMessage struct {
    Message

    didAck bool
    ack    func()
}

func NewMessage(data interface{}, attributes ...string) *SubMessage {
    return nil
}

func NewStringMessage(data string, attributes ...string) *SubMessage {
    return &SubMessage{}
}

func (m *SubMessage) SetAck(ack func()) *SubMessage {
    if m != nil {
        m.ack = ack
    }
    return m
}

func (m *SubMessage) Ack() {
    if m != nil && !m.didAck {
        m.ack()
        m.didAck = true
    }
}

func (m *SubMessage) SetAttribute(key string, value string) *SubMessage {
    if m != nil {
        if m.Attributes == nil {
            m.Attributes = make(map[string]*core.Value)
        }
        m.Attributes[key] = core.NewStringValue(value)
    }
    return m
}
