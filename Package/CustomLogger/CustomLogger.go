package CustomLogger

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"

	"go.uber.org/zap"
)

type CustomLoggerInterface interface {
	CustomLog()
}

type CustomLoggerStruct struct {
}

func NewLogger(Mode int) (*zap.Logger, error) {

	switch Mode {
	case DevMode.Client:
		devOpts := zap.Development()
		logger, err := zap.NewDevelopment(devOpts)
		if err != nil {
			return nil, err
		}
		var zapLogger = zap.Must(logger, err)
		if len(err.Error()) > 1 {
			return nil, err
		}
		return zapLogger, nil
	case DevMode.QA:
		return nil, nil
	case DevMode.PROD:
		return nil, nil
	}
	return nil, nil

}
