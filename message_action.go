package main

import (
	"fmt"
	"github.com/nlopes/slack"
)

type Message struct {
	Type string `json:"type"`
	User string `json:"user"`
	Ts   string `json:"ts"`
	Text string `json:"text"`
}

type MessageAction struct {
	Token       string               `json:"token"`
	CallbackID  string               `json:"callback_id"`
	Type        string               `json:"type"`
	TriggerID   string               `json:"trigger_id"`
	ResponseURL string               `json:"response_url"`
	Team        MessageActionTeam    `json:"team"`
	Channel     MessageActionChannel `json:"channel"`
	User        MessageActionUser    `json:"user"`
	Message     Message              `json:"message"`
}

type ReportGroupDataTimestamp struct {
	ReportGroupData
	Timestamp string
	Params    *PostMessageParameters
}

var groupReportMessages = make(map[string]ReportGroupDataTimestamp)

func (message *MessageAction) Handle(rtm *slack.RTM) error {
	switch callbackID := message.CallbackID; callbackID {
	case SpamCallback:
		reported := message.Message.User
		data := ReportGroupData{
			Reporter: message.User.ID,
			Reported: reported,
			Channel:  message.Channel,
			HasSeen:  false,
			Removed:  false,
		}
		params, timestamp, err := SendTemplate(
			rtm,
			ReportGroup.ID,
			ReportGroupTemplate,
			data,
		)
		if err != nil {
			ReportError(rtm, message.ResponseURL)
			return err
		}
		groupReportMessages[reported] = ReportGroupDataTimestamp{
			ReportGroupData: data,
			Timestamp:       timestamp,
			Params:          params,
		}

		_, _, err = DMTemplate(
			rtm,
			reported,
			ReportIMTemplate,
			ReportIMData{
				Channel: message.Channel,
			},
		)
		if err != nil {
			ReportError(rtm, message.ResponseURL)
			return err
		}

		actionResponse := &MessageActionResponse{
			Text:         "The user has anonymously been notified of your report. The moderator team has also received your report, thank you very much for the report.",
			ResponseType: "ephemeral",
		}
		return actionResponse.Send(message.ResponseURL)
	}
	return fmt.Errorf("received unknown callback ID: %s", message.CallbackID)
}

// ReportError reports the error to the user if it can.
func ReportError(rtm *slack.RTM, responseURL string) error {
	actionResponse := &MessageActionResponse{
		Text:         "Could not process the report, please notify the administrators manually or try again.",
		ResponseType: Ephemeral,
	}
	return actionResponse.Send(responseURL)
}

// UpdateGroupReport updates the report posted in the report group
func UpdateGroupReport(rtm *slack.RTM, reported string, removed bool) error {
	params, err := TemplateToMessage(
		ReportGroupTemplate,
		ReportGroupData{
			HasSeen:  true,
			Removed:  removed,
			Reported: reported,
		},
	)
	if err != nil {
		return err
	}

	data := groupReportMessages[reported]

	_, _, _, err = rtm.Client.SendMessage(
		ReportGroup.ID,
		slack.MsgOptionUpdate(data.Timestamp),
		slack.MsgOptionText(data.Params.Text, data.Params.EscapeText),
		slack.MsgOptionAttachments(params.Attachments...),
	)
	return err
}
