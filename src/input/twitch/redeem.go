package twitch

import (
	"encoding/json"
	"strings"
	"time"
	"twitch-redeem-trigger/src/config"
	"twitch-redeem-trigger/src/logger"

	"github.com/gorilla/websocket"
	"github.com/nicklaw5/helix/v2"
)

func ListenForRedeems(cfgTwitch config.Twitch, l *logger.Logger, events chan<- RedeemEvent) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     cfgTwitch.Auth.ClientID,
		ClientSecret: cfgTwitch.Auth.ClientSecret,
	})
	if err != nil {
		l.Error("Failed to create Twitch client: " + err.Error())
		return
	}

	// Automatisches Token-Handling
	token := strings.TrimPrefix(cfgTwitch.Auth.OAuth, "oauth:")

	if cfgTwitch.Auth.RefreshToken != "" {
		if cfgTwitch.Auth.ClientSecret == "" {
			l.Error("TWITCH_REFRESH_TOKEN is set, but TWITCH_CLIENT_SECRET is missing. Refreshing requires the Secret.")
		} else {
			l.Info("Refreshing User Access Token from Twitch using Refresh Token...")
			resp, err := client.RefreshUserAccessToken(cfgTwitch.Auth.RefreshToken)
			if err != nil {
				l.Error("Failed to refresh token: %v", err)
			} else {
				token = resp.Data.AccessToken
				l.Info("Successfully refreshed User Access Token.")
			}
		}
	} else if (token == "" || cfgTwitch.Auth.ClientSecret != "") && cfgTwitch.Auth.RefreshToken == "" {
		l.Info("Requesting App Access Token from Twitch...")
		resp, err := client.RequestAppAccessToken([]string{})
		if err != nil {
			l.Error("Failed to request App Access Token: %v", err)
		} else {
			token = resp.Data.AccessToken
			l.Info("Successfully obtained App Access Token (Note: Might not work for WebSockets).")
		}
	}

	if token == "" {
		l.Error("No valid OAuth token available. Please provide TWITCH_OAUTH_TOKEN or a valid Refresh Token + Secret.")
		return
	}
	client.SetUserAccessToken(token)

	// Resolve Broadcaster ID if missing
	if cfgTwitch.Channel.Id == "" {
		l.Error("Channel ID missing, resolving for name: %s", cfgTwitch.Channel.Name)
		return
	}

	if cfgTwitch.Channel.Id == "" {
		l.Error("Broadcaster ID is required for EventSub subscription. Please check your .env file (TWITCH_USER_ID).")
		return
	}

	wsURL := "wss://eventsub.wss.twitch.tv/ws"

	for {
		l.Info("Connecting to Twitch EventSub WebSocket: %s", wsURL)
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			l.Error("Failed to connect to WebSocket: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			wsURL = "wss://eventsub.wss.twitch.tv/ws" // Reset URL on failure
			continue
		}

		sessionID := ""

		// Message handling loop
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				l.Error("WebSocket read error: %v", err)
				break
			}

			var msg websocketMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				l.Error("Failed to unmarshal WebSocket message: %v", err)
				continue
			}

			switch msg.Metadata.MessageType {
			case "session_welcome":
				var payload sessionPayload
				json.Unmarshal(msg.Payload, &payload)
				sessionID = payload.Session.ID
				l.Info("Connected to EventSub WebSocket. SessionID: %s", sessionID)

				// Subscribe to Channel Points Redeems
				resp, err := client.CreateEventSubSubscription(&helix.EventSubSubscription{
					Type:    helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd,
					Version: "1",
					Condition: helix.EventSubCondition{
						BroadcasterUserID: cfgTwitch.Channel.Id,
					},
					Transport: helix.EventSubTransport{
						Method:    "websocket",
						SessionID: sessionID,
					},
				})
				if err != nil {
					l.Error("Failed to create EventSub subscription request: %v", err)
				} else if resp.StatusCode >= 400 {
					l.Error("Failed to create EventSub subscription: %s (Status: %d, Error: %s)", resp.ErrorMessage, resp.StatusCode, resp.Error)
					if resp.StatusCode == 401 {
						l.Error("Authentication failed. This often happens if the TWITCH_CLIENT_ID and TWITCH_OAUTH_TOKEN do not match.")
						l.Error("If you used twitchtokengenerator.com, ensure you use their Client ID (%s) or create a token with your own Client ID.", cfgTwitch.Auth.ClientID)
					}
				} else {
					l.Info("Subscribed to channel points redeems for channel %s (ID: %s)", cfgTwitch.Channel.Name, cfgTwitch.Channel.Id)
				}

			case "notification":
				var payload notificationPayload
				json.Unmarshal(msg.Payload, &payload)

				if payload.Subscription.Type == helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd {
					var event helix.EventSubChannelPointsCustomRewardRedemptionEvent
					json.Unmarshal(payload.Event, &event)

					l.Debug("Redeem received: %s by %s", event.Reward.Title, event.UserName)

					if event.Reward.Title == cfgTwitch.Redeem.Name {
						l.Info("Target redeem triggered: %s", event.Reward.Title)
						events <- RedeemEvent{
							RedeemName: event.Reward.Title,
							User:       event.UserName,
						}
					}
				}

			case "session_reconnect":
				var payload sessionPayload
				json.Unmarshal(msg.Payload, &payload)
				wsURL = payload.Session.ReconnectURL
				l.Info("Reconnect requested. New URL: %s", wsURL)
				conn.Close()
				goto nextConn

			case "session_keepalive":
				l.Debug("WebSocket keepalive received")
			}
		}
		conn.Close()
	nextConn:
		time.Sleep(1 * time.Second)
	}
}
