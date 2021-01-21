package triggers

import (
	mevt "maunium.net/go/mautrix/event"
)

// pattern stolen from https://github.com/whyrusleeping/hellabot/blob/master/hellabot.go
type SyncTrigger struct {
	// Returns true if this trigger applies to the passed in message
	Cond func(*mevt.Event) (shouldApply bool)

	// The action to perform if Cond is true
	// return true if processing should continue
	Act func(*mevt.Event) (shouldApply bool, response *mevt.MessageEventContent)

	Meta *TriggerMeta
}

func (t *SyncTrigger) Condition(msg *mevt.Event) (shouldApply bool) {
	return t.Cond(msg)
}
func (t *SyncTrigger) Action(msg *mevt.Event) (shouldApply bool, response *mevt.MessageEventContent) {
	return t.Act(msg)
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
func (t *ComposedTrigger) Condition(msg *mevt.Event) (shouldApply bool) {
	for i := range t.subTriggers {
		if t.subTriggers[i].Condition(msg) {
			return true
		}
	}
	return false
}
func (t *ComposedTrigger) Action(msg *mevt.Event) (shouldApply bool, response *mevt.MessageEventContent) {
	for i := range t.subTriggers {
		if t.subTriggers[i].Condition(msg) {
			return t.subTriggers[i].Action(msg)
		}
	}
	return true, nil // not possible
}
func (t *ComposedTrigger) GetMeta() *TriggerMeta {
	return t.meta
}
func (t *ComposedTrigger) SetMeta(m *TriggerMeta) {
	t.meta = m
}


type Trigger interface {
	Condition(msg *mevt.Event) (shouldApply bool)
	Action(msg *mevt.Event) (shouldApply bool, response *mevt.MessageEventContent)
	GetMeta() *TriggerMeta
	SetMeta(*TriggerMeta)
}


type TriggerMeta struct {
	Disabled bool
	Name string
}