package main

import (
	"encoding/json"
	"io/ioutil"
)

// BotConfig required when starting the bot.
type BotConfig struct {
	VerificationToken string `json:"verification_token"`
	BotUserOAuthToken string `json:"bot_user_oauth_token"`
	ReportGroupName   string `json:"report_group_name"`

	Port int `json:"port"`
}

// Config currently in use.
var Config BotConfig

// SetupConfig needs to be called to initialize the config.
func SetupConfig() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}
