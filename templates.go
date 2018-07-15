package main

import (
	"text/template"
)

var ReportGroupTemplate *template.Template
var ReportIMTemplate *template.Template
var RulesAccepted *template.Template
var WelcomeTemplate *template.Template

type ReportIMData struct {
	Channel MessageActionChannel
}

type ReportGroupData struct {
	Reporter string
	Reported string
	Channel  MessageActionChannel
	HasSeen  bool
	Removed  bool
}

// SetupTemplates is required to be called to parse the templates.
func SetupTemplates() {
	ReportGroupTemplate = template.Must(
		template.ParseFiles("templates/report-group.json"),
	)
	ReportIMTemplate = template.Must(
		template.ParseFiles("templates/report-im.json"),
	)
	RulesAccepted = template.Must(
		template.ParseFiles("templates/rules-accepted.json"),
	)
	WelcomeTemplate = template.Must(
		template.ParseFiles("templates/welcome.json"),
	)
}
