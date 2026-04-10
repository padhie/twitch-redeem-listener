package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func Load(envFile string) (*Config, error) {
	// Falls kein Pfad angegeben wurde, nutze Standard .env
	if envFile == "" {
		envFile = ".env"
	}

	// Lade angegebene env-Datei
	if err := godotenv.Load(envFile); err != nil {
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
				IP:          os.Getenv("TAPO_IP"),
				Username:    os.Getenv("TAPO_USERNAME"),
				Password:    os.Getenv("TAPO_PASSWORD"),
				TriggerWord: os.Getenv("TAPO_TRIGGER_WORD"),
			},
			Media: Media{
				Port:     getEnvInt("MEDIA_PORT", 0),
				Mappings: parseMappings(),
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
			Address: os.Getenv("WEB_BIND_ADDRESS"),
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

func parseMappings() map[string]string {
	mappings := make(map[string]string)

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "MEDIA_REWARD_") {
			continue
		}

		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		rewardName := strings.TrimPrefix(key, "MEDIA_REWARD_")
		rewardName = strings.ReplaceAll(rewardName, "_", " ")

		mappings[rewardName] = value
	}

	return mappings
}
