package CustomLogger

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"
)

type CustomLoggerInterface interface {
	CustomLog()
}

type CustomLoggerStruct struct {
}

func NewLogger(Mode int) (CustomLoggerStruct, error) {
	logger := CustomLoggerStruct{}
	switch Mode {
	case DevMode.Client:
		/*
			devOpts := zap.Development()
			logger, err := zap.NewDevelopment(devOpts)
			if err != nil {
				return nil, err
			}
			var zapLogger = zap.Must(logger, err)
			if len(err.Error()) > 1 {
				return nil, err
			}*/
		return logger, nil
	case DevMode.QA:
		return logger, nil
	case DevMode.PROD:
		return logger, nil
	}
	return logger, nil

}
func (log CustomLoggerStruct) CustomLog() {}
