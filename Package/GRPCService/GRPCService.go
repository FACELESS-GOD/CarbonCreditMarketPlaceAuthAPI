package GRPCService

import (
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"
)

type ControllerStruct struct {
	Mdl  Model.ModelStruct
	Conf Configurator.ConfiguratorStruct
	Log  CustomLogger.CustomLoggerInterface
}

func NewCtrl(Mdl Model.ModelStruct, Conf Configurator.ConfiguratorStruct, Log CustomLogger.CustomLoggerInterface) ControllerStruct {
	ctrl := ControllerStruct{}
	ctrl.Conf = Conf
	ctrl.Mdl = Mdl
	ctrl.Log = Log
	return ctrl
}
