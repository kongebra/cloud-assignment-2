package main

import (
	"fmt"
	"log"
	"net/http"
)

type Webhook struct {
	ID string `json:"-"`
	WebhookURL string `json:"webhookURL"`
	MinTriggerValue int `json:"minTriggerValue"`
	Trigger int `json:"-"`
}

func (w *Webhook) CheckTriggerValue() {
	if w.MinTriggerValue <= 0 {
		w.MinTriggerValue = 1
	}

	w.Trigger = 0
}

func (w *Webhook) CheckTrigger() bool {
	w.Trigger++

	if w.Trigger >= w.MinTriggerValue {
		w.Trigger = 0
		return true
	}

	return false
}

func (w *Webhook) SendHook() {
	fmt.Println(w.WebhookURL)

	response, err := http.Get(w.WebhookURL)


	if err != nil {
		log.Fatal("Something went wrong sending a hook-request")
	}

	defer response.Body.Close()
}