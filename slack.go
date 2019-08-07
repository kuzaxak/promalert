package main

import (
	"bytes"
	"github.com/nlopes/slack"
	"strings"
	"text/template"
)

func SlackSendAlertMessage(alert Alert, token, channel, messageTemplate string, images ...SlackImage) (string, string, error) {
	api := slack.New(token)

	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
	}

	t, err := template.New("message").Funcs(funcMap).Parse(messageTemplate)
	fatal(err, "error in template")
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, alert); err != nil {
		return "", "", err
	}

	textBlockObj := slack.NewTextBlockObject(
		"mrkdwn",
		tpl.String(),
		false,
		false,
	)

	var blocks []slack.Block

	blocks = append(blocks, slack.NewDividerBlock())
	blocks = append(blocks, slack.NewSectionBlock(textBlockObj, nil, nil))

	for _, image := range images {
		textBlock := slack.NewTextBlockObject("plain_text", image.Title, false, false)
		blocks = append(blocks, slack.NewImageBlock(image.Url, "metric graph "+image.Title, "", textBlock))
	}

	messageBlocks := slack.MsgOptionBlocks(blocks...)
	respChannel, respTimestamp, err := api.PostMessage(channel, messageBlocks)
	return respChannel, respTimestamp, err
}
