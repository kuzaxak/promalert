package main

import (
	"fmt"
	"github.com/mitchellh/hashstructure"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strconv"
	"time"
)

func (alert Alert) Hash() string {
	hash, err := hashstructure.Hash(map[string]KV{
		"labels":      alert.Labels,
		"annotations": alert.Annotations,
	}, nil)
	fatal(err, "Hash cant be calculated")

	return strconv.FormatUint(hash, 10)
}

func (alert Alert) GeneratePictures() []SlackImage {
	generatorUrl, err := url.Parse(alert.GeneratorURL)
	if err != nil {
		panic(err)
	}

	generatorQuery, _ := url.ParseQuery(generatorUrl.RawQuery)
	var alertFormula string
	for key, param := range generatorQuery {
		if key == "g0.expr" {
			alertFormula = param[0]
			break
		}
	}
	fmt.Println(alertFormula)

	plotExpression := GetPlotExpr(alertFormula)
	queryTime, duration := alert.GetPlotTimeRange()

	var images []SlackImage

	for _, expr := range plotExpression {
		plot := Plot(
			expr,
			queryTime,
			duration,
			time.Duration(viper.GetInt64("metric_resolution")),
			viper.GetString("prometheus_url"),
			alert,
		)

		publicURL, err := UploadFile(viper.GetString("s3_bucket"), viper.GetString("s3_region"), plot)
		fatal(err, "failed to upload")
		log.Printf("Graph uploaded, URL: %s", publicURL)

		images = append(images, SlackImage{
			Url:   publicURL,
			Title: expr.String(),
		})
	}

	return images
}

func (alert Alert) PostMessage() (string, string, []slack.Block) {
	log.Printf("Alert: status=%s,Labels=%v,Annotations=%v", alert.Status, alert.Labels, alert.Annotations)
	severity := alert.Labels["severity"]
	log.Printf("no action on severity: %s", severity)

	options := make([]slack.MsgOption, 0)

	if alert.Status == AlertStatusFiring || alert.MessageTS == "" {
		messageBlocks, err := ComposeMessageBody(
			alert,
			viper.GetString("message_template"),
			viper.GetString("header_template"),
			alert.GeneratePictures()...,
		)
		fatal(err, "failed to generate slack message")

		alert.MessageBody = messageBlocks
		options = append(options, slack.MsgOptionBlocks(messageBlocks...))
		if alert.MessageTS != "" {
			options = append(options, slack.MsgOptionBroadcast())
		}
	} else {
		messageBody, err := ComposeResolveUpdateBody(
			alert,
			viper.GetString("header_template"),
			alert.GeneratePictures()...,
		)
		fatal(err, "failed to generate slack message")
		options = append(options, messageBody)
	}

	if alert.MessageTS != "" {
		options = append(options, slack.MsgOptionTS(alert.MessageTS))

		updateBlocks := alert.MessageBody
		d, err := ComposeUpdateFooter(alert, viper.GetString("footer_template"))
		fatal(err, "failed to generate slack message")
		updateBlocks = append(updateBlocks, d...)

		respChannel, respTimestamp, err := SlackUpdateAlertMessage(
			viper.GetString("slack_token"),
			alert.Channel,
			alert.MessageTS,
			slack.MsgOptionBlocks(updateBlocks...),
		)
		fatal(err, "failed to send slack message")
		log.Printf("Slack message updated, channel: %s thread: %s", respChannel, respTimestamp)
	}

	channel := viper.GetString("slack_channel")

	if alert.Channel != "" {
		channel = alert.Channel
	}

	respChannel, respTimestamp, err := SlackSendAlertMessage(
		viper.GetString("slack_token"),
		channel,
		options...,
	)
	fatal(err, "failed to send slack message")
	log.Printf("Slack message sended, channel: %s thread: %s", respChannel, respTimestamp)

	return respChannel, respTimestamp, alert.MessageBody
}

func (alert Alert) GetPlotTimeRange() (time.Time, time.Duration) {
	var queryTime time.Time
	var duration time.Duration
	if alert.StartsAt.Second() > alert.EndsAt.Second() {
		queryTime = alert.StartsAt
		duration = time.Minute * 20
	} else {
		queryTime = alert.EndsAt
		duration = queryTime.Sub(alert.StartsAt)

		if duration < time.Minute*20 {
			duration = time.Minute * 20
		}
	}
	log.Printf("Querying Time %v Duration: %v", queryTime, duration)
	return queryTime, duration
}
