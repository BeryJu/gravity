package extconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func RecoverWrapper(input any) any {
	if input != nil {
		go ReportEatenCat()
	}
	return input
}

func ReportEatenCat() {
	log.Println("Reporting eaten cat")
	data := map[string]string{
		"domain": "gravity-panic",
		"name":   "pageview",
		"url":    "http://localhost",
	}
	d, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "https://analytics.beryju.org/api/event", bytes.NewBuffer(d))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("gravity %s", FullVersion()))
	http.DefaultClient.Do(req)
}
