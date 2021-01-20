package ping

import (
	"github.com/matrix-org/go-neb/clients"
	"github.com/matrix-org/go-neb/services/utils"
	mevt "maunium.net/go/mautrix/event"
	"strings"
)

var userPingPong = &utils.SyncTrigger{
	Cond: func(b *clients.BotClient, event *mevt.Event) (shouldApply bool) {
		msg := event.Content.AsMessage()
		return msg.MsgType == mevt.MsgText && strings.HasPrefix(msg.Body, "!ping")
	},
	Act: func(botClient *clients.BotClient, event *mevt.Event) (shouldApply bool) {
		botClient.SendMessageEvent(event.RoomID, mevt.EventMessage, &mevt.MessageEventContent{
			MsgType: mevt.MsgNotice,
			Body:    "Pong!",
		})
		return false
	},
	Meta: &utils.TriggerMeta{
		Disabled: false,
		Name:     "userPingPong",
	},
}

