package main

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// StartServer starts the server listening to webhooks.
func StartServer(rtm *slack.RTM) {
	http.HandleFunc("/interactive", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Error(err)
			return
		}
		payload := r.FormValue("payload")
		eventsAPIEvent, err := slackevents.ParseEvent(
			json.RawMessage(payload),
			slackevents.OptionVerifyToken(
				&slackevents.TokenComparator{
					VerificationToken: Config.VerificationToken,
				},
			),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Error(err)
			return
		}

		if eventsAPIEvent.Type == InteractiveMessageType {
			var interactiveMessage InteractiveMessage
			err = json.Unmarshal([]byte(payload), &interactiveMessage)
			if err != nil {
				logrus.Error(err)
				return
			}
			err = interactiveMessage.Handle(rtm)
			if err != nil {
				logrus.Error(err)
				return
			}
		} else if eventsAPIEvent.Type == MessageActionType {
			var messageAction MessageAction
			err = json.Unmarshal([]byte(payload), &messageAction)
			if err != nil {
				logrus.Error(err)
				return
			}

			err := messageAction.Handle(rtm)
			if err != nil {
				logrus.Error(err)
				return
			}

		} else {
			logrus.Warnf("Received unknown API event type: %s\n", eventsAPIEvent.Type)
		}
	})
	logrus.Info("Server listening")
	http.ListenAndServe(":"+strconv.Itoa(Config.Port), nil)
}
