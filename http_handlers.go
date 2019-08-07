package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Ok!")
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if viper.GetBool("debug") {
		// Save a copy of this request for debugging.
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		}
		log.Printf("New request")
		fmt.Println(string(requestDump))
	}

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var m HookMessage
	if err := dec.Decode(&m); err != nil {
		log.Printf("error decoding message: %v", err)
		http.Error(w, "invalid request body", 400)
		return
	}

	log.Printf("Alerts: GroupLabels=%v, CommonLabels=%v", m.GroupLabels, m.CommonLabels)

	for _, alert := range m.Alerts {
		log.Printf("Alert: status=%s,Labels=%v,Annotations=%v", alert.Status, alert.Labels, alert.Annotations)
		severity := alert.Labels["severity"]
		log.Printf("no action on severity: %s", severity)

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
		queryTime, duration := GetPlotTimeRange(alert)

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

		respChannel, respTimestamp, err := SlackSendAlertMessage(
			alert,
			viper.GetString("slack_token"),
			viper.GetString("slack_channel"),
			viper.GetString("message_template"),
			images...,
		)
		fatal(err, "failed to send slack message")
		log.Printf("Slack message sended, channel: %s thread: %s", respChannel, respTimestamp)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("{\"success\": true}"))
	fatal(err, "failed to send response")
}
