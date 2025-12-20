package Router

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"
	"CarbonCreditMarketPlaceAuthAPI/Package/Controller"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"

	"github.com/gin-gonic/gin"
)

func NewRouter(Mdl Model.ModelStruct, Mode int) (*gin.Engine, error) {
	ctrl := Controller.ControllerStruct{}
	ctrl.Mdl = Mdl
	router := gin.Default()

	router.Use(gin.Recovery())

	router.GET("/Verify", ctrl.EditUserCred)

	authorized := router.Group("/")

	if Mode != DevMode.Client {

		authorized.Use(ctrl.AuthMiddleware())
		{

			router.POST("/Add", ctrl.AddUser)
			router.DELETE("/Delete", ctrl.DeleteUser)
			router.PUT("/EditUsr", ctrl.EditUser)
			router.PUT("/EditUsrPass", ctrl.EditUserCred)

		}

	} else {

		// For Testing

		router.POST("/Add", ctrl.AddUser)
		router.DELETE("/Delete", ctrl.DeleteUser)
		router.PUT("/EditUsr", ctrl.EditUser)
		router.PUT("/EditUsrPass", ctrl.EditUserCred)

	}

	return router, nil
}
