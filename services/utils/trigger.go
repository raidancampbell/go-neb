package utils

import (
	"github.com/matrix-org/go-neb/clients"
	mevt "maunium.net/go/mautrix/event"
)

// pattern stolen from https://github.com/whyrusleeping/hellabot/blob/master/hellabot.go
type SyncTrigger struct {
	// Returns true if this trigger applies to the passed in message
	Cond func(*clients.BotClient, *mevt.Event) (shouldApply bool)

	// The action to perform if Cond is true
	// return true if processing should continue
	Act func(*clients.BotClient, *mevt.Event) (shouldApply bool)

	Meta *TriggerMeta
}

func (t *SyncTrigger) Condition(b *clients.BotClient, msg *mevt.Event) (shouldApply bool) {
	return t.Cond(b, msg)
}
func (t *SyncTrigger) Action(b *clients.BotClient, msg *mevt.Event) (shouldApply bool) {
	return t.Act(b, msg)
}
func (t *SyncTrigger) GetMeta() *TriggerMeta {
	return t.Meta
}
func (t *SyncTrigger) SetMeta(m *TriggerMeta) {
	t.Meta = m
}


type ComposedTrigger struct {
	subTriggers []Trigger
	meta *TriggerMeta
}
func (t *ComposedTrigger) Condition(b *clients.BotClient, msg *mevt.Event) (shouldApply bool) {
	for i := range t.subTriggers {
		if t.subTriggers[i].Condition(b, msg) {
			return true
		}
	}
	return false
}
func (t *ComposedTrigger) Action(b *clients.BotClient, msg *mevt.Event) (shouldApply bool) {
	for i := range t.subTriggers {
		if t.subTriggers[i].Condition(b, msg) {
			return t.subTriggers[i].Action(b, msg)
		}
	}
	return true // not possible
}
func (t *ComposedTrigger) GetMeta() *TriggerMeta {
	return t.meta
}
func (t *ComposedTrigger) SetMeta(m *TriggerMeta) {
	t.meta = m
}


type Trigger interface {
	Condition(b *clients.BotClient, msg *mevt.Event) (shouldApply bool)
	Action(b *clients.BotClient, msg *mevt.Event) (shouldApply bool)
	GetMeta() *TriggerMeta
	SetMeta(*TriggerMeta)
}


type TriggerMeta struct {
	Disabled bool
	Name string
}