package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	SetupConfig()
	SetupTemplates()

	api := slack.New(Config.BotUserOAuthToken)

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()

	go rtm.ManageConnection()
	go WatchEvents(rtm)
	StartServer(rtm)
}
