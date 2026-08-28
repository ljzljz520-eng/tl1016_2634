package notification

import (
	"fmt"
	"labops/internal/model"
	"sort"
	"strings"
)

type Message struct {
	ID, Recipient, Subject, Body, Channel, State string
	Attempts                                     int
}
type Outbox struct{ messages map[string]Message }

func New() *Outbox { return &Outbox{messages: map[string]Message{}} }
func (o *Outbox) Queue(m Message) error {
	if m.ID == "" || m.Recipient == "" {
		return fmt.Errorf("message identity required")
	}
	if m.Subject == "" {
		return fmt.Errorf("subject required")
	}
	if m.Channel == "" {
		m.Channel = "inbox"
	}
	if _, ok := o.messages[m.ID]; ok {
		return fmt.Errorf("message exists")
	}
	m.State = "queued"
	o.messages[m.ID] = m
	return nil
}
func (o *Outbox) Get(id string) (Message, error) {
	m, ok := o.messages[id]
	if !ok {
		return Message{}, fmt.Errorf("message missing")
	}
	return m, nil
}
func (o *Outbox) Deliver(id string) error {
	m, e := o.Get(id)
	if e != nil {
		return e
	}
	if m.State != "queued" && m.State != "retry" {
		return fmt.Errorf("message not deliverable")
	}
	m.State = "sent"
	m.Attempts++
	o.messages[id] = m
	return nil
}
func (o *Outbox) Retry(id string) error {
	m, e := o.Get(id)
	if e != nil {
		return e
	}
	if m.State != "queued" && m.State != "failed" {
		return fmt.Errorf("message cannot retry")
	}
	m.State = "retry"
	m.Attempts++
	o.messages[id] = m
	return nil
}
func (o *Outbox) Fail(id string) error {
	m, e := o.Get(id)
	if e != nil {
		return e
	}
	if m.State == "sent" {
		return fmt.Errorf("message sent")
	}
	m.State = "failed"
	o.messages[id] = m
	return nil
}
func (o *Outbox) Pending() []Message {
	out := []Message{}
	for _, m := range o.messages {
		if m.State == "queued" || m.State == "retry" {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func (o *Outbox) ForRecipient(recipient string) []Message {
	out := []Message{}
	for _, m := range o.messages {
		if m.Recipient == recipient {
			out = append(out, m)
		}
	}
	return out
}
func (o *Outbox) Count(state string) int {
	n := 0
	for _, m := range o.messages {
		if state == "" || m.State == state {
			n++
		}
	}
	return n
}
func Render(r model.Record, action string) Message {
	return Message{ID: r.ID + "-" + action, Recipient: r.Owner, Subject: action + " " + r.AssetTag, Body: fmt.Sprintf("Record %s is %s", r.ID, action), Channel: "inbox"}
}
func NormalizeChannel(c string) string {
	c = strings.ToLower(strings.TrimSpace(c))
	switch c {
	case "email", "sms", "inbox":
		return c
	}
	return "inbox"
}
