package Router

import (
	"CarbonCreditMarketPlaceAuthAPI/Package/Controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(Ctrl Controller.ControllerStruct, Mode int) (*gin.Engine, error) {

	router := gin.Default()

	router.Use(gin.Recovery())

	router.GET("/Verify", Ctrl.VerifyCred)

	authorized := router.Group("/")

	authorized.Use(Ctrl.AuthMiddleware())
	{

		router.POST("/Add", Ctrl.AddUser)
		router.DELETE("/Delete", Ctrl.DeleteUser)
		router.PUT("/EditUsr", Ctrl.EditUser)
		router.PUT("/EditUsrPass", Ctrl.EditUserCred)
	}

	return router, nil
}
