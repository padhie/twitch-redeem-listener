package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	// Lade .env-Datei
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		Twitch: Twitch{
			Channel: TwitchChannel{
				Id:   os.Getenv("TWITCH_CHANNEL_ID"),
				Name: os.Getenv("TWITCH_CHANNEL_NAME"),
			},
			Auth: TwitchAuth{
				ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
				ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
				OAuth:        os.Getenv("TWITCH_OAUTH_TOKEN"),
				RefreshToken: os.Getenv("TWITCH_REFRESH_TOKEN"),
			},
			Redeem: TwitchRedeem{
				Name:   os.Getenv("REDEEM_NAME"),
				Status: os.Getenv("REDEEM_STATUS"),
			},
		},

		Output: Output{
			Type: os.Getenv("OUTPUT_TYPE"),
			Tapo: Tapo{
				IP:       os.Getenv("TAPO_IP"),
				Username: os.Getenv("TAPO_USERNAME"),
				Password: os.Getenv("TAPO_PASSWORD"),
			},
		},

		GPIO: GPIO{
			Enabled:   getEnvBool("GPIO_ENABLE"),
			RedeemPin: getEnvInt("LED_REDEEM_GPIO", 17),
			TapoPin:   getEnvInt("LED_TAPO_GPIO", 18),
		},

		Logging: Logging{
			LogFile:   os.Getenv("LOG_FILE"),
			UseSyslog: getEnvBool("USE_SYSLOG"),
		},

		Web: Web{
			Enabled: getEnvBool("ENABLE_WEB_INTERFACE"),
			Port:    os.Getenv("WEB_PORT"),
		},
	}

	return cfg, nil
}

func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}

	return value
}

func getEnvBool(key string) bool {
	value := os.Getenv(key)

	switch value {
	case "true", "1", "yes", "y", "on":
		return true
	default:
		return false
	}
}
