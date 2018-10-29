package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const DiscordWebhookUrl = "https://discordapp.com/api/webhooks/506567057899520000/ZJ0BuT5z0cTFRzUULiPqXDBNvgh4mL1mwDdTgd2Gy9mT_fOl46agHEKjs5qvp6J8lK55"

type WebhookInfo struct {
	Content string `json:"content"`
}

func SendDiscordLogEntry(message string) {
	info := WebhookInfo{}
	info.Content = message + "\n"
	raw, _ := json.Marshal(info)

	response, err := http.Post(DiscordWebhookUrl, "application/json", bytes.NewBuffer(raw))

	if err != nil {
		fmt.Println(err)
		fmt.Println(ioutil.ReadAll(response.Body))
	}
}