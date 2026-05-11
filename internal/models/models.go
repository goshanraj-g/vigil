package models

import "time"

type Monitor struct {
	ID         string `json:"id"`
	Query      string `json:"query"`
	WebhookURL     string `json:"webhook_url"`
	Schedule       string `json:"schedule"`
	Status         string `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type Result struct {
	ID            string `json:"id"`
	MonitorID     string `json:"monitor_id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	Snippet       string `json:"snippet"`
	FoundAt       time.Time `json:"foundAt"`
}