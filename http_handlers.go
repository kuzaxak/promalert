package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Ok!")
}

func webhook(w http.ResponseWriter, r *http.Request) {
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
		var alertLevel float64
		var alertOperator string

		for key, param := range generatorQuery {
			if key == "g0.expr" {
				alertFormula = param[0]
				break
			}
		}

		fmt.Println(alertFormula)
		expr, _ := promql.ParseExpr(alertFormula)

		if binaryExpr, ok := expr.(*promql.BinaryExpr); ok {
			alertFormula = binaryExpr.LHS.String()
			alertLevel, _ = strconv.ParseFloat(binaryExpr.RHS.String(), 64)

			if binaryExpr.Op == promql.ItemLTE || binaryExpr.Op == promql.ItemLSS {
				alertOperator = "LE"
			} else {
				alertOperator = "GE"
			}
		}

		// Fetch from Prometheus
		log.Printf("Querying Prometheus %s", alertFormula)

		var queryTime time.Time
		var duration time.Duration

		if alert.StartsAt.Second() > alert.EndsAt.Second() {
			queryTime = alert.StartsAt
			duration = time.Minute * 10
		} else {
			queryTime = alert.EndsAt
			duration = queryTime.Sub(alert.StartsAt)

			if duration < time.Minute*10 {
				duration = time.Minute * 10
			}
		}

		metrics, err := Metrics(
			viper.GetString("prometheus_url"),
			alertFormula,
			queryTime,
			duration,
			time.Duration(viper.GetInt64("metric_resolution")),
		)
		fatal(err, "failed to get metrics")

		var selectedMetrics model.Matrix

		for _, metric := range metrics {
			founded := false
			for label, value := range metric.Metric {
				if originValue, ok := alert.Labels[string(label)]; ok {
					if originValue == string(value) {
						founded = true
					} else {
						founded = false
						break
					}
				}
			}

			if founded {
				selectedMetrics = model.Matrix{metric}
				break
			}
		}

		// Plot
		log.Printf("Creating plot %s", alert.Annotations["summary"])
		plot, err := Plot(selectedMetrics, alertLevel, alertOperator)
		fatal(err, "failed to create plot")

		publicURL, err := UploadFile(viper.GetString("s3_bucket"), viper.GetString("s3_region"), plot)

		_, _, err = SlackSendAlertMessage(
			alert,
			viper.GetString("slack_token"),
			viper.GetString("slack_channel"),
			publicURL,
			viper.GetString("message_template"),
		)
		fatal(err, "failed to send slack message")
	}

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("{\"success\": true}"))
	fatal(err, "failed to send response")
}
