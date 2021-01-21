package ping

import (
	"github.com/matrix-org/go-neb/services/utils/triggers"
	mevt "maunium.net/go/mautrix/event"
	"strings"
)

var UserPingPong = &triggers.SyncTrigger{
	Cond: func(event *mevt.Event) (shouldApply bool) {
		msg := event.Content.AsMessage()
		return msg.MsgType == mevt.MsgText && strings.HasPrefix(msg.Body, "!ping")
	},
	Act: func(event *mevt.Event) (shouldApply bool, response *mevt.MessageEventContent) {
		return false, &mevt.MessageEventContent{
			MsgType: mevt.MsgNotice,
			Body:    "Pong!",
		}
	},
	Meta: &triggers.TriggerMeta{
		Disabled: false,
		Name:     "userPingPong",
	},
}

