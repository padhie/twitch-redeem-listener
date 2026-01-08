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

func (d NoopDevice) Toggle() error {
	d.logger.Info("Noop: Toggle")

	return nil
}
