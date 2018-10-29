package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Discord Webhook URL
const DiscordWebhookUrl = "https://discordapp.com/api/webhooks/506567057899520000/ZJ0BuT5z0cTFRzUULiPqXDBNvgh4mL1mwDdTgd2Gy9mT_fOl46agHEKjs5qvp6J8lK55"

// Webhook Information stuct
type WebhookInfo struct {
	Content string `json:"content"`
}

// Sends discord log entry
func SendDiscordLogEntry(message string) {
	// Create WebhookInfo
	info := WebhookInfo{}
	// Sets the content
	info.Content = message + "\n"
	// Encode JSON
	raw, _ := json.Marshal(info)

	// Post the JSON to URL
	response, err := http.Post(DiscordWebhookUrl, "application/json", bytes.NewBuffer(raw))

	// Check if an error
	if err != nil {
		// Print error
		fmt.Println(err)
		// Print response body
		fmt.Println(ioutil.ReadAll(response.Body))
	}
}