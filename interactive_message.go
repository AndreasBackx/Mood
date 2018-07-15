package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nlopes/slack"
)

type InteractiveMessage struct {
	Type             string                   `json:"type"`
	Actions          []slack.AttachmentAction `json:"actions"`
	CallbackID       string                   `json:"callback_id"`
	Team             MessageActionTeam        `json:"team"`
	Channel          MessageActionChannel     `json:"channel"`
	User             MessageActionUser        `json:"user"`
	ActionTimestamp  json.Number              `json:"action_ts"`
	MessageTimestamp json.Number              `json:"message_ts"`
	AttachmentID     json.Number              `json:"attachment_id"`
	Token            string                   `json:"token"`
	OriginalMessage  slack.Message            `json:"original_message"`
	ResponseURL      string                   `json:"response_url"`
	TriggerID        string                   `json:"trigger_id"`
}

func (message *InteractiveMessage) Handle(rtm *slack.RTM) error {
	if len(message.Actions) == 0 {
		return errors.New("No actions passed in interactive message")
	}

	action := message.Actions[0]

	switch callbackID := message.CallbackID; callbackID {
	case RulesCallback:
		if action.Value == "yes" {
			DMTemplate(
				rtm,
				message.User.ID,
				RulesAccepted,
				struct{}{},
			)
		}

	case SpamResponseCallback:
		var actionResponse *MessageActionResponse
		removed := action.Value == "removed"
		reported := message.User.ID

		err := UpdateGroupReport(rtm, reported, removed)
		if err != nil {
			return err
		}

		if removed {
			actionResponse = &MessageActionResponse{
				Text:            "Thank you for removing your crossposts, I've forwarded it to the moderator team.",
				ReplaceOriginal: true,
			}
		} else if action.Value == "invalid" {
			actionResponse = &MessageActionResponse{
				Text:            "Thank you for marking the report has invalid, I've forwarded it to the moderator team.",
				ReplaceOriginal: true,
			}
		}
		return actionResponse.Send(message.ResponseURL)
	}
	return fmt.Errorf("received unknown callback ID: %s", message.CallbackID)
}
