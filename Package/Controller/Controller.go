package Controller

import (
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddUserResponseStruct struct {
	UserID int `json:"userid"`
}

type AddUserRequestStruct struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

func (Ctrl *ControllerStruct) AddUser(gCtx *gin.Context) {

	var req AddUserRequestStruct = AddUserRequestStruct{}

	err := gCtx.ShouldBindJSON(&req)

	if err != nil {
		gCtx.JSON(http.StatusBadRequest, req)
		return
	}

	mdlReq := Model.ModelAddUserRequestStruct{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	res := Ctrl.Mdl.AddUser(mdlReq)

	if Ctrl.Mdl.IsAnyError != true {
		gCtx.JSON(http.StatusInternalServerError, req)
		return
	}

	if res.UserID < 1 {
		gCtx.JSON(http.StatusInternalServerError, req)
		return
	}

	tkn, err := Ctrl.Mdl.AddToken(res.UserID)

	if err != nil {
		gCtx.JSON(http.StatusInternalServerError, req)
		return
	}

	gCtx.Header("Authorization", "Bearer "+tkn)

	gCtx.JSON(http.StatusOK, AddUserResponseStruct{UserID: res.UserID})

}
