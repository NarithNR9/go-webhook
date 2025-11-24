package webhook

// UQPayEvent represents the structure of a UQPay webhook event.
type UQPayEvent struct {
	Version   string                 `json:"version"`
	EventType string                 `json:"event_type"`
	EventName string                 `json:"event_name"`
	EventID   string                 `json:"event_id"`
	Data      map[string]interface{} `json:"data"`
}
