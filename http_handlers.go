package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
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
		if prevAlert, founded := FindAlert(alert); founded {
			alert.Channel = prevAlert.Channel
			alert.MessageTS = prevAlert.MessageTS
			alert.MessageBody = prevAlert.MessageBody
			respChannel, respTimestamp, _ := alert.PostMessage()
			if alert.Status == AlertStatusFiring {
				alert.MessageTS = respTimestamp
				alert.Channel = respChannel
				AddAlert(alert)
			}
			log.Printf("Slack update sended, channel: %s thread: %s", respChannel, respTimestamp)
		} else {
			// post new message
			respChannel, respTimestamp, messageBody := alert.PostMessage()
			alert.MessageTS = respTimestamp
			alert.Channel = respChannel
			alert.MessageBody = messageBody

			AddAlert(alert)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("{\"success\": true}"))
	fatal(err, "failed to send response")
}
