package output

import (
	"strings"

	"twitch-redeem-trigger/src/config"
	"twitch-redeem-trigger/src/logger"
)

func Build(cfgOutput config.Output, l *logger.Logger) Device {
	var outputType = strings.ToUpper(cfgOutput.Type)
	l.Debug("Output type: %s", outputType)

	if outputType == "TAPO" {
		return BuildTapo(cfgOutput.Tapo, l)
	}

	if outputType == "MEDIA" {
		return BuildMedia(cfgOutput.Media, l)
	}

	return BuildNoop(l)
}
