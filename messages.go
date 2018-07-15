package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"text/template"
)

// PostMessageParameters reoresenting a message sent to a user.
type PostMessageParameters struct {
	slack.PostMessageParameters

	Text string `json:"text"`
}

type MessageActionTeam struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type MessageActionUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type MessageActionChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const (
	RulesCallback        = "rules"
	SpamCallback         = "spam"
	SpamResponseCallback = "spam_response"

	MessageActionType      = "message_action"
	InteractiveMessageType = "interactive_message"
)

// TemplateToMessage converts a template to a message.
func TemplateToMessage(tmpl *template.Template, data interface{}) (*PostMessageParameters, error) {
	if tmpl == nil {
		return nil, fmt.Errorf("Template is nil, did you call SetupTemplates?")
	}

	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, data)
	if err != nil {
		return nil, err
	}

	var params PostMessageParameters
	err = json.Unmarshal(buffer.Bytes(), &params)
	if err != nil {
		return nil, err
	}

	return &params, nil
}

// DMTemplate sends one direct message to a user.
func DMTemplate(
	rtm *slack.RTM,
	userID string,
	tmpl *template.Template,
	data interface{},
) (*PostMessageParameters, string, error) {

	_, _, imChannelID, err := rtm.OpenIMChannel(userID)
	if err != nil {
		return nil, "", err
	}

	params, timestamp, err := SendTemplate(rtm, imChannelID, tmpl, data)
	if err != nil {
		return params, timestamp, err
	}

	_, _, err = rtm.CloseIMChannel(imChannelID)
	if err != nil {
		return params, timestamp, err
	}

	return params, timestamp, nil
}

// SendTemplate parses a template and sends it to a channel.
// Returns the timestamp of the sent message and an optional error.
func SendTemplate(
	rtm *slack.RTM,
	channelID string,
	tmpl *template.Template,
	data interface{},
) (*PostMessageParameters, string, error) {
	params, err := TemplateToMessage(tmpl, data)
	if err != nil {
		return nil, "", err
	}

	_, timestamp, err := rtm.PostMessage(channelID, params.Text, params.PostMessageParameters)
	if err != nil {
		return nil, "", err
	}
	return params, timestamp, nil
}
