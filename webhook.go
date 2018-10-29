package main

type Webhook struct {
	WebhookURL string `json:"webhookURL"`
	MinTriggerValue int `json:"minTriggerValue"`
}

func (w *Webhook) CheckTriggerValue() {
	if w.MinTriggerValue <= 0 {
		w.MinTriggerValue = 1
	}
}