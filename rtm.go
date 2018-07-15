package main

import (
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

var SlackInfo *slack.Info
var TestGroup slack.Group
var ReportGroup slack.Group

func WatchEvents(rtm *slack.RTM) {

	for msg := range rtm.IncomingEvents {
		// fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {

		case *slack.ConnectedEvent:
			SlackInfo = ev.Info
			groups, _ := rtm.GetGroups(true)

			testGroupFound, reportGroupFound := false, false
			for _, group := range groups {
				if group.Name == "mood-test" {
					TestGroup = group
					testGroupFound = true
				}
				if group.Name == Config.ReportGroupName {
					ReportGroup = group
					reportGroupFound = true
				}
			}
			if !testGroupFound {
				logrus.Error("Test group not found")
				return
			}
			if !reportGroupFound {
				logrus.Error("Report group not found")
				return
			}

		case *slack.MessageEvent:
			if ev.Channel == TestGroup.ID {
				user, err := rtm.GetUserInfo(ev.User)
				if err != nil {
					logrus.Error(err)
					continue
				}
				_, _, err = DMTemplate(
					rtm,
					ev.User,
					WelcomeTemplate,
					struct {
						User *slack.User
						Team *slack.Team
					}{
						User: user,
						Team: SlackInfo.Team,
					},
				)
				if err != nil {
					logrus.Error(err)
					continue
				}
			} else {
				logrus.Infof("Message not in test group: %s != %s", ev.Channel, TestGroup.ID)
			}

		case *slack.InvalidAuthEvent:
			logrus.Error("Invalid credentials")
			return
		}
	}
}
