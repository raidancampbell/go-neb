package htmltitle

import (
	"bytes"
	"github.com/matrix-org/go-neb/services/utils/triggers"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	mevt "maunium.net/go/mautrix/event"
	"net/http"
	"net/url"
	"strings"
)

const (
	maxHTMLResponseBytes = 1024 * 1024 * 5 // 5 MB
)

// prints the html title text of any URLs within a message
var HTMLTitle = &triggers.SyncTrigger{
	Cond: func(event *mevt.Event) bool {
		msg := event.Content.AsMessage()
		return msg.MsgType == mevt.MsgText && extractURLWord(msg.Body) != ""
	},

	Act: func(event *mevt.Event) (shouldContinue bool, response *mevt.MessageEventContent) {
		msg := event.Content.AsMessage().Body
		shouldContinue = true // always
		var u *url.URL
		for _, word := range strings.Split(msg, " ") {
			if strings.Contains(word, "://"){
				tmpURL, err := url.Parse(strings.Trim(word, ":,!.<>"))
				if err == nil {
					u = tmpURL
					break
				}
			}
		}
		// not possible
		if u == nil {
			return
		}
		logrus.Debugf("retrieving title for URL '%s'", u.String())

		r, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			logrus.WithError(err).Errorf("failed to create HTTP request to endpoint %+v", u)
			return
		}

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			logrus.WithError(err).Errorf("failed to complete HTTP request to endpoint %+v", u)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, maxHTMLResponseBytes))
		if err != nil {
			logrus.WithError(err).Errorf("failed to read %d bytes from %+v", maxHTMLResponseBytes, u)
			return
		}
		logrus.Debugf("URL %s had an HTML body of length %d", u.String(), len(body))
		tok := html.NewTokenizer(bytes.NewReader(body))
		for {
			tokType := tok.Next()
			if tokType == html.ErrorToken {
				if tok.Err() != io.EOF {
					logrus.WithError(tok.Err()).Error("error while parsing HTML response")
				}
				break
			}

			t := tok.Token()
			if t.Data == "title" {
				nextType := tok.Next()
				if nextType == html.TextToken {
					return false, &mevt.MessageEventContent{
						MsgType: mevt.MsgNotice,
						Body:    strings.Trim(tok.Token().Data, "\r\n\t "),
					}
				}
			}
		}
		return
	},
	Meta: &triggers.TriggerMeta{
		Disabled: false,
		Name:     "HTMLTitle",
	},
}

func extractURLWord(s string) string {
	for _, word := range strings.Split(s, " ") {
		if strings.Contains(word, "://"){
			_, err := url.Parse(strings.Trim(word, ":,!.<>"))
			if err == nil {
				return word
			}
		}
	}
	return ""
}