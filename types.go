package main

import "time"

type HookMessage struct {
	Version           string  `json:"version"`
	GroupKey          string  `json:"groupKey"`
	Status            string  `json:"status"`
	Receiver          string  `json:"receiver"`
	GroupLabels       KV      `json:"groupLabels"`
	CommonLabels      KV      `json:"commonLabels"`
	CommonAnnotations KV      `json:"commonAnnotations"`
	ExternalURL       string  `json:"externalURL"`
	Alerts            []Alert `json:"alerts"`
}

// Alert holds one alert for notification templates.
type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
}

// KV is a set of key/value string pairs.
type KV map[string]string
