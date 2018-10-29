package main

import (
	"log"
	"net/http"
)

// Webhook Struct
type Webhook struct {
	ID string `json:"-"`
	WebhookURL string `json:"webhookURL"`
	MinTriggerValue int `json:"minTriggerValue"`
	Trigger int `json:"-"`
}

// Check the default trigger value
func (w *Webhook) CheckTriggerValue() {
	if w.MinTriggerValue <= 0 {
		w.MinTriggerValue = 1
	}

	w.Trigger = 0
}

// Check the current trigger
func (w *Webhook) CheckTrigger() bool {
	w.Trigger++

	// Trigger has reached it's limit
	if w.Trigger >= w.MinTriggerValue {
		w.Trigger = 0
		return true
	}

	return false
}

// Send the webhook
func (w *Webhook) SendHook() {
	// get response from the request
	response, err := http.Get(w.WebhookURL)

	// Check for errors
	if err != nil {
		log.Fatal("Something went wrong sending a hook-request")
	}

	// Close up
	defer response.Body.Close()
}