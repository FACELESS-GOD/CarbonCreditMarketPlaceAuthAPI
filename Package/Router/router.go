package Router

import (
	"CarbonCreditMarketPlaceAuthAPI/Package/Controller"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"

	"github.com/gin-gonic/gin"
)

func NewRouter(Mdl Model.ModelStruct) {
	ctrl := Controller.ControllerStruct{}
	ctrl.Mdl = Mdl
	router := gin.Default()
	router.POST("/Add", ctrl.AddUser)

}
