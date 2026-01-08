package gpio

import (
	"fmt"
	"time"

	"twitch-redeem-trigger/config"
	"twitch-redeem-trigger/logger"

	"github.com/warthog618/go-gpiocdev"
)

const (
	RedeemLED = "redeem"
	TapoLED   = "tapo"
)

var (
	enabled   bool
	redeemPin *gpiocdev.Line
	tapoPin   *gpiocdev.Line
	chip      *gpiocdev.Chip
	log       *logger.Logger
)

func Init(cfgGPIO config.GPIO, l *logger.Logger) error {
	log = l
	enabled = cfgGPIO.Enabled
	if !enabled {
		return nil
	}

	l.Info("GPIO enabled (Raspberry Pi 5 compatible mode)")

	// Auf dem Pi 5 ist das meistens /dev/gpiochip4 für die GPIO-Leiste
	// Wir versuchen den Chip zu öffnen
	var err error
	chip, err = gpiocdev.NewChip("/dev/gpiochip4")
	if err != nil {
		// Fallback auf gpiochip0, falls gpiochip4 nicht existiert
		chip, err = gpiocdev.NewChip("/dev/gpiochip0")
		if err != nil {
			return fmt.Errorf("failed to open gpiochip: %v", err)
		}
	}

	// Konfiguriere Pins
	redeemPin, err = chip.RequestLine(cfgGPIO.RedeemPin, gpiocdev.AsOutput())
	if err != nil {
		return fmt.Errorf("failed to request redeem pin %d: %v", cfgGPIO.RedeemPin, err)
	}
	redeemPin.SetValue(0)

	tapoPin, err = chip.RequestLine(cfgGPIO.TapoPin, gpiocdev.AsOutput())
	if err != nil {
		return fmt.Errorf("failed to request tapo pin %d: %v", cfgGPIO.TapoPin, err)
	}
	tapoPin.SetValue(0)

	return nil
}

func Close() {
	if !enabled {
		return
	}

	if redeemPin != nil {
		redeemPin.SetValue(0)
		redeemPin.Close()
	}
	if tapoPin != nil {
		tapoPin.SetValue(0)
		tapoPin.Close()
	}
	if chip != nil {
		chip.Close()
	}
}

func BlinkLED(ledType string, duration time.Duration) {
	if !enabled {
		return
	}

	var pin *gpiocdev.Line
	switch ledType {
	case RedeemLED:
		pin = redeemPin
	case TapoLED:
		pin = tapoPin
	default:
		log.Error("Unknown LED type: %s", ledType)
		return
	}

	if pin == nil {
		return
	}

	// Blinke die LED
	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		pin.SetValue(1)
		time.Sleep(200 * time.Millisecond)
		pin.SetValue(0)
		time.Sleep(200 * time.Millisecond)
	}
}
