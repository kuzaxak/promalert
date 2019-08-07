package main

import (
	"bytes"
	"github.com/nlopes/slack"
	"strings"
	"text/template"
)

func SlackSendAlertMessage(alert Alert, token, channel, publicURL, message_template string) (string, string, error) {
	api := slack.New(token)

	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
	}

	t, err := template.New("message").Funcs(funcMap).Parse(message_template)
	fatal(err, "error in template")
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, alert); err != nil {
		return "", "", err
	}

	textBlock := slack.NewTextBlockObject("plain_text", alert.Annotations["summary"], false, false)
	textBlockObj := slack.NewTextBlockObject(
		"mrkdwn",
		tpl.String(),
		false,
		false,
	)

	messageBlocks := slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		slack.NewSectionBlock(textBlockObj, nil, nil),
		slack.NewImageBlock(publicURL, "metric graph", "", textBlock),
	)

	respChannel, respTimestamp, err := api.PostMessage(channel, messageBlocks)
	return respChannel, respTimestamp, err
}

//
//func SlackSendToThread(token) error {
//	api := slack.New(token)
//	_ = api
//
//	return nil
//}
//
//
