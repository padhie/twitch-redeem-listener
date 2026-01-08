package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"twitch-redeem-trigger/src/config"
	"twitch-redeem-trigger/src/logger"
)

type Tapo struct {
	IP       string
	Username string
	Password string
	Client   *http.Client
	logger   *logger.Logger
}

func BuildTapo(cfgTapo config.Tapo, l *logger.Logger) Device {
	return &Tapo{
		IP:       cfgTapo.IP,
		Username: cfgTapo.Username,
		Password: cfgTapo.Password,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: l,
	}
}

func (d *Tapo) Toggle() error {
	// 1. Anmeldung (vereinfacht)
	loginURL := fmt.Sprintf("http://%s", d.IP)
	loginPayload := map[string]interface{}{
		"method": "login",
		"params": map[string]string{
			"username": d.Username,
			"password": d.Password,
		},
	}

	loginReq, err := json.Marshal(loginPayload)
	if err != nil {
		return err
	}

	resp, err := d.Client.Post(loginURL, "application/json", bytes.NewBuffer(loginReq))
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("login returned status: %d", resp.StatusCode)
	}

	// 2. Steckdose toggle (vereinfacht)
	toggleURL := fmt.Sprintf("http://%s", d.IP)
	togglePayload := map[string]interface{}{
		"method": "passthrough",
		"params": map[string]interface{}{
			"device_id": d.IP,
			"requestData": map[string]interface{}{
				"system": map[string]interface{}{
					"set_relay_state": map[string]interface{}{
						"state": 1, // 1 = an, 0 = aus
					},
				},
			},
		},
	}

	toggleReq, err := json.Marshal(togglePayload)
	if err != nil {
		return err
	}

	resp, err = d.Client.Post(toggleURL, "application/json", bytes.NewBuffer(toggleReq))
	if err != nil {
		return fmt.Errorf("toggle failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("toggle returned status: %d", resp.StatusCode)
	}

	// 3. 10 Sekunden warten und ausschalten
	time.Sleep(10 * time.Second)

	togglePayload["params"].(map[string]interface{})["requestData"] = map[string]interface{}{
		"system": map[string]interface{}{
			"set_relay_state": map[string]interface{}{
				"state": 0, // aus
			},
		},
	}

	toggleReq, err = json.Marshal(togglePayload)
	if err != nil {
		return err
	}

	resp, err = d.Client.Post(toggleURL, "application/json", bytes.NewBuffer(toggleReq))
	if err != nil {
		return fmt.Errorf("toggle off failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("toggle off returned status: %d", resp.StatusCode)
	}

	d.logger.Info("Tapo device toggled successfully")
	return nil
}
