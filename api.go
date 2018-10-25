package main

type API struct {
	Uptime string `json:"uptime"`
	Info string `json:"info"`
	Version string `json:"version"`
}

func (api *API) CalculateUptime(startTime int) {
	api.Uptime = ConvertSecondsToISO8601(startTime)
}