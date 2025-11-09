package messaging

import (
	"github.com/mojo-lang/mojo/go/pkg/mojo/core"
)

type SubMessage struct {
	Message

	didAck     bool
	ack        func()
	nck        func()
	term       func()
	inProgress func()
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

func (m *SubMessage) SetNak(nak func()) *SubMessage {
	if m != nil {
		m.nck = nak
	}
	return m
}

func (m *SubMessage) SetTerm(term func()) *SubMessage {
	if m != nil {
		m.term = term
	}
	return m
}

func (m *SubMessage) SetInProgress(inProgress func()) *SubMessage {
	if m != nil {
		m.inProgress = inProgress
	}
	return m
}

func (m *SubMessage) Ack() {
	if m != nil && !m.didAck && m.ack != nil {
		m.ack()
		m.didAck = true
	}
}

func (m *SubMessage) Nak() {
	if m != nil && !m.didAck && m.nck != nil {
		m.nck()
		m.didAck = true
	}
}

func (m *SubMessage) Term() {
	if m != nil && !m.didAck && m.term != nil {
		m.term()
		m.didAck = true
	}
}

func (m *SubMessage) InProgress() {
	if m != nil && m.inProgress != nil {
		m.inProgress()
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
