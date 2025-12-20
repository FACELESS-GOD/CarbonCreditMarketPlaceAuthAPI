package main

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"
	"CarbonCreditMarketPlaceAuthAPI/Package/Router"
	"log"
)

func main() {
	//Mode := DevMode.QA
	//Mode := DevMode.PROD
	Mode := DevMode.Client

	conf, err := Configurator.NewConfiguration(Mode)

	if err != nil {
		log.Panic(err.Error())
	}

	mdl := Model.NewModel(conf, CustomLogger.CustomLoggerStruct{})

	router, err := Router.NewRouter(mdl, Mode)

	if err != nil {
		log.Panic(err.Error())
	}

	if err := router.Run(conf.ADDRESS); err != nil {
		log.Fatal(err)
	}
}
