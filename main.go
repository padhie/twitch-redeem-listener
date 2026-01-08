package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"twitch-redeem-trigger/config"
	"twitch-redeem-trigger/gpio"
	"twitch-redeem-trigger/input/twitch"
	"twitch-redeem-trigger/logger"
	"twitch-redeem-trigger/output"
	"twitch-redeem-trigger/web"
)

func main() {
	// 1. Konfiguration laden
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Logger initialisieren
	l := logger.Build(cfg.Logging)
	l.Info("Starting Twitch Redeem Trigger...")

	// 3. GPIO/LEDs initialisieren (vielleicht auslagern)
	if err := gpio.Init(cfg.GPIO, l); err != nil {
		l.Error("Failed to init GPIO: %v. Continuing without LEDs.", err)
	}
	defer gpio.Close()

	initGpioTest(l)

	// 4. output initialisieren
	device := output.Build(cfg.Output, l)

	// 5. Web-Server starten (non-blocking)
	initWebServer(cfg.Web, device, l)

	// 6. Twitch-Redeem-Listener starten
	initTwitch(cfg.Twitch, l, device)

	// 7. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	l.Info("Shutting down...")
}

func initGpioTest(l *logger.Logger) {
	l.Info("Starting GPIO test...")

	l.Info("RedeemLED blink now.")
	go gpio.BlinkLED(gpio.RedeemLED, 3*time.Second)

	l.Info("TapoLED blink now.")
	go gpio.BlinkLED(gpio.TapoLED, 3*time.Second)

	l.Info("GPIO test finished.")
}

func initWebServer(cfgWeb config.Web, device output.Device, l *logger.Logger) {
	if cfgWeb.Enabled {
		go func() {
			if err := web.StartServer(cfgWeb, device, l); err != nil {
				l.Error("Web server failed: %v", err)
			}
		}()
	}
}

func initTwitch(cfgTwitch config.Twitch, l *logger.Logger, device output.Device) {
	redeemChan := make(chan twitch.RedeemEvent)
	go twitch.ListenForRedeems(cfgTwitch, l, redeemChan)

	for event := range redeemChan {
		l.Info("Received redeem: %s (User: %s)", event.RedeemName, event.User)

		// LED für Redeem blinken lassen
		go gpio.BlinkLED(gpio.RedeemLED, 3*time.Second)

		// Tapo-Steckdose schalten
		if err := device.Toggle(); err != nil {
			l.Error("Failed to toggle Tapo device: %v", err)
			// Neu starten (wie besprochen)
			os.Exit(1)
		}

		// LED für Tapo blinken lassen
		go gpio.BlinkLED(gpio.TapoLED, 3*time.Second)
	}
}
