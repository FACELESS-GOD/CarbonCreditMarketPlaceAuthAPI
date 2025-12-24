package main

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/Controller"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"
	"CarbonCreditMarketPlaceAuthAPI/Package/Router"
	"log"
)

func main() {
	//Mode := DevMode.QA
	//Mode := DevMode.PROD
	//Mode := DevMode.Client
	Mode := DevMode.Test

	EnvPath := "./"

	conf, err := Configurator.NewConfiguration(Mode, EnvPath)

	if err != nil {
		log.Panic(err.Error())
	}

	mdl := Model.NewModel(conf, CustomLogger.CustomLoggerStruct{})

	ctrl := Controller.ControllerStruct{}
	ctrl.Mdl = mdl

	router, err := Router.NewRouter(ctrl, Mode)

	if err != nil {

		
		log.Panic(err.Error())
	}

	if err := router.Run(conf.ADDRESS); err != nil {
		log.Fatal(err)
	}

}
