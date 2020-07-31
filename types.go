package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"time"
)

type HookMessage struct {
	Version           string  `json:"version"`
	GroupKey          string  `json:"groupKey"`
	Status            string  `json:"status" binding:"required"`
	Receiver          string  `json:"receiver"`
	GroupLabels       KV      `json:"groupLabels"`
	CommonLabels      KV      `json:"commonLabels"`
	CommonAnnotations KV      `json:"commonAnnotations"`
	ExternalURL       string  `json:"externalURL"`
	Alerts            []Alert `json:"alerts" binding:"required"`
}

// Alert holds one alert for notification templates.
type Alert struct {
	Status       AlertStatus `json:"status" binding:"required"`
	Labels       KV          `json:"labels"`
	Annotations  KV          `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt" binding:"required"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL" binding:"required"`
	Fingerprint  string      `json:"fingerprint"`
	Channel      string
	MessageTS    string
	MessageBody  []slack.Block
}

type AlertStatus string

const (
	AlertStatusFiring   AlertStatus = "firing"
	AlertStatusResolved AlertStatus = "resolved"
)

// KV is a set of key/value string pairs.
type KV map[string]string

type SlackImage struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

type PlotExpr struct {
	Formula  string
	Operator string
	Level    float64
}

func (expr PlotExpr) String() string {
	return fmt.Sprintf("%s %s %.2f", expr.Formula, expr.Operator, expr.Level)
}
