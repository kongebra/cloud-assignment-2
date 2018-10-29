package main

type WebhookURL struct {
	Type string `json:"type"`
}

type MinTriggerValue struct {
	Type int `json:"type"`
}

type Webhook struct {
	WebhookURL WebhookURL `json:"webhookURL"`
	MinTriggerValue MinTriggerValue `json:"minTriggerValue"`
}