package worker

import "encoding/json"

const TypeSearchMonitor = "monitor:search"

type SearchMonitorPayload struct {
	MonitorID string `json:"monitor_id"`
}

func NewSearchPayload(monitorID string) ([]byte, error) {
	return json.Marshal(SearchMonitorPayload{MonitorID: monitorID})
}
