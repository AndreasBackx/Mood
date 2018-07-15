package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MessageActionResponse struct {
	Text            string             `json:"text,omitempty"`
	Attachments     []slack.Attachment `json:"attachments,omitempty"`
	ThreadTs        string             `json:"thread_ts,omitempty"`
	ResponseType    string             `json:"response_type"`
	ReplaceOriginal bool               `json:"replace_original"`
	DeleteOriginal  bool               `json:"delete_original"`
}

const (
	InChannel = "in_channel"
	Ephemeral = "ephemeral"
)

var client = &http.Client{}

func (actionResponse *MessageActionResponse) Send(responseURL string) error {
	logrus.Infof("responseURL: %s", responseURL)
	body, err := json.Marshal(actionResponse)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		responseURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		logrus.Warn(response)
		return fmt.Errorf("Slack API returned a status of %d", response.StatusCode)
	}

	return nil
}
