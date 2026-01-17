package output

import (
	"twitch-redeem-trigger/src/logger"
)

type NoopDevice struct {
	logger *logger.Logger
}

func BuildNoop(l *logger.Logger) Device {
	return NoopDevice{
		logger: l,
	}
}

func (d NoopDevice) Toggle(input ToggleInput) error {
	d.logger.Info("Noop: Toggle")
	d.logger.Info("User: %s", input.User)
	d.logger.Info("Redeem: %s", input.RedeemName)

	return nil
}
