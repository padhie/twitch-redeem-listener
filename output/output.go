package output

import (
	"strings"

	"twitch-redeem-trigger/config"
	"twitch-redeem-trigger/logger"
)

func Build(cfgOutput config.Output, l *logger.Logger) Device {
	var outputType = strings.ToUpper(cfgOutput.Type)

	if outputType == "TAPO" {
		return BuildTapo(cfgOutput.Tapo, l)
	}

	return BuildNoop(l)
}
