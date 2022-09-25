package extconfig

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Recover() any {
	err := recover()
	if err == nil {
		return err
	}
	_, ok := err.(error)
	if !ok {
		return err
	}
	go ReportEatenCat()
	return nil
}

func ReportEatenCat() {
	data := map[string]string{
		"domain":   "gravity-panic",
		"name":     "pageview",
		"referrer": FullVersion(),
	}
	d, err := json.Marshal(data)
	if err != nil {
		return
	}
	http.DefaultClient.Post("https://analytics.beryju.org/api/event", "application/json", bytes.NewBuffer(d))
}
