package main

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/slack-go/slack"
	"github.com/spf13/cast"
)

func SlackSendAlertMessage(token, channel string, messageOptions ...slack.MsgOption) (string, string, error) {
	api := slack.New(token)
	respChannel, respTimestamp, err := api.PostMessage(channel, messageOptions...)
	return respChannel, respTimestamp, err
}

func SlackUpdateAlertMessage(token, channel, timestamp string, messageOptions ...slack.MsgOption) (string, string, error) {
	api := slack.New(token)
	respChannel, respTimestamp, respText, err := api.UpdateMessage(channel, timestamp, messageOptions...)

	_ = respText

	return respChannel, respTimestamp, err
}

func ComposeResolveUpdateBody(alert Alert, headerTemplate string, images ...SlackImage) (slack.MsgOption, error) {
	headerTpl, e := ParseTemplate(headerTemplate, alert)
	if e != nil {
		return nil, e
	}
	statusBlock := slack.NewTextBlockObject(
		"mrkdwn",
		headerTpl.String(),
		false,
		false,
	)

	var blocks []slack.Block
	blocks = append(blocks, slack.NewSectionBlock(statusBlock, nil, nil))
	blocks = append(blocks, slack.NewDividerBlock())
	for _, image := range images {
		textBlock := slack.NewTextBlockObject("plain_text", image.Title, false, false)
		blocks = append(blocks, slack.NewImageBlock(image.Url, "metric graph "+image.Title, "", textBlock))
	}
	messageBlocks := slack.MsgOptionBlocks(blocks...)

	return messageBlocks, nil
}

func ComposeUpdateFooter(alert Alert, footerTemplate string) ([]slack.Block, error) {
	footerTpl, e := ParseTemplate(footerTemplate, alert)
	if e != nil {
		return nil, e
	}
	footerBlock := slack.NewTextBlockObject(
		"mrkdwn",
		footerTpl.String(),
		false,
		false,
	)

	var blocks []slack.Block
	blocks = append(blocks, slack.NewDividerBlock())
	blocks = append(blocks, slack.NewContextBlock("", footerBlock))

	return blocks, nil
}

func ComposeMessageBody(alert Alert, messageTemplate, headerTemplate string, images ...SlackImage) ([]slack.Block, error) {
	tpl, e := ParseTemplate(messageTemplate, alert)
	if e != nil {
		return nil, e
	}
	headerTpl, e := ParseTemplate(headerTemplate, alert)
	if e != nil {
		return nil, e
	}
	statusBlock := slack.NewTextBlockObject(
		"mrkdwn",
		headerTpl.String(),
		false,
		false,
	)

	textBlockObj := slack.NewTextBlockObject(
		"mrkdwn",
		tpl.String(),
		false,
		false,
	)
	var blocks []slack.Block
	blocks = append(blocks, slack.NewSectionBlock(statusBlock, nil, nil))
	blocks = append(blocks, slack.NewDividerBlock())
	blocks = append(blocks, slack.NewSectionBlock(textBlockObj, nil, nil))
	for _, image := range images {
		textBlock := slack.NewTextBlockObject("plain_text", image.Title, false, false)
		blocks = append(blocks, slack.NewImageBlock(image.Url, "metric graph "+image.Title, "", textBlock))
	}

	return blocks, nil
}

func ParseTemplate(messageTemplate string, alert Alert) (bytes.Buffer, error) {
	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
		"now":     time.Now,
		"dateFormat": func(layout string, v interface{}) (string, error) {
			t, err := cast.ToTimeE(v)
			if err != nil {
				return "", err
			}

			return t.Format(layout), nil
		},
	}
	msgTpl, err := template.New("message").Funcs(funcMap).Parse(messageTemplate)
	fatal(err, "error in template")
	var tpl bytes.Buffer
	if err := msgTpl.Execute(&tpl, alert); err != nil {
		return bytes.Buffer{}, err
	}

	return tpl, err
}
