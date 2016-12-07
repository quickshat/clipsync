package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type activeDevice struct {
	IP       string
	Port     int64
	LastPing time.Time
	Metrics  metrics
}

func (a *activeDevice) loadMetrics() error {
	resp, err := http.Get("http://" + a.IP + ":" + fmt.Sprint(a.Port) + "/metrics")
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(b, &a.Metrics)
	fmt.Println(a.Metrics)
	return nil
}

func (a *activeDevice) register() error {
	emitLog("WEB", "Register on new device")
	r, err := http.PostForm("http://"+a.IP+":"+fmt.Sprint(a.Port)+"/register", url.Values{
		"port": {"8081"},
	})
	if err != nil || r.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
