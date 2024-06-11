package zaplogger

import (
	"github.com/mohsenHa/messenger/logger/loggerentity"
	"go.uber.org/zap"
)

func logParamsToZapParams(keys map[loggerentity.ExtraKey]interface{}) []zap.Field {
	params := make([]zap.Field, 0, len(keys))

	for k, v := range keys {
		params = append(params, zap.Any(string(k), v))
	}

	return params
}
