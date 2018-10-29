package main

// API stuct
type API struct {
	Uptime string `json:"uptime"`
	Info string `json:"info"`
	Version string `json:"version"`
}

// Calculate the uptime of the application
func (api *API) CalculateUptime(startTime int) {
	api.Uptime = ConvertSecondsToISO8601(startTime)
}