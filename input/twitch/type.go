package twitch

import (
	"encoding/json"
	"time"

	"github.com/nicklaw5/helix/v2"
)

type RedeemEvent struct {
	RedeemName string
	User       string
}

// Twitch WebSocket Message Structures
type metadata struct {
	MessageID        string    `json:"message_id"`
	MessageType      string    `json:"message_type"`
	MessageTimestamp time.Time `json:"message_timestamp"`
}

type sessionPayload struct {
	Session struct {
		ID                      string    `json:"id"`
		Status                  string    `json:"status"`
		ConnectedAt             time.Time `json:"connected_at"`
		KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
		ReconnectURL            string    `json:"reconnect_url"`
	} `json:"session"`
}

type notificationPayload struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Event        json.RawMessage            `json:"event"`
}

type websocketMessage struct {
	Metadata metadata        `json:"metadata"`
	Payload  json.RawMessage `json:"payload"`
}
