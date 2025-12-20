package Controller

import (
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AddUserResponseStruct struct {
	UserID int `json:"userid"`
}

type AddUserRequestStruct struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DeleteUserRequestStruct struct {
	UserID int `json:"userid" binding:"required"`
}

type DeleteUserResponseStruct struct {
	IsDeleted     bool
	IsAnyError    bool
	ErrorMessages []string
}

type EditUserRequestStruct struct {
	UserID            int    `json:"userid" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Email             string `json:"email" binding:"required"`
	Password          string `json:"password" binding:"required"`
	Ispasswordchanged bool   `json:"ispasswordchanged"`
}

type EditUserResponseStruct struct {
	UserDt        EditUserRequestStruct
	IsDeleted     bool
	IsAnyError    bool
	ErrorMessages []string
}

type VerifyUserRequestStruct struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUserResponseStruct struct {
	IsDeleted     bool
	IsAnyError    bool
	ErrorMessages []string
}

type ControllerStruct struct {
	Mdl  Model.ModelStruct
	Conf Configurator.ConfiguratorStruct
	Log  CustomLogger.CustomLoggerInterface
}

/*
DeleteUser
	EditUser(ModelEditUserRequestStruct) error
	UpdateCred(ModelUpdateCredRequestStruct) error
	VerifyToken(Token string, UserID int) (bool, error)
	UpdateToken(UserID int, Token string) (bool, error)
	VerifyCred(ModelVerifyCredRequestStruct) (bool, error)
*/

func NewCtrl(Mdl Model.ModelStruct, Conf Configurator.ConfiguratorStruct, Log CustomLogger.CustomLoggerInterface) ControllerStruct {
	ctrl := ControllerStruct{}
	ctrl.Conf = Conf
	ctrl.Mdl = Mdl
	ctrl.Log = Log
	return ctrl
}

func (Ctrl *ControllerStruct) AddUser(gCtx *gin.Context) {

	var req AddUserRequestStruct = AddUserRequestStruct{}

	err := gCtx.BindJSON(&req)

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

	if Ctrl.Mdl.IsAnyError == true {
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

func (Ctrl *ControllerStruct) DeleteUser(gCtx *gin.Context) {

	var req DeleteUserRequestStruct = DeleteUserRequestStruct{}

	var res DeleteUserResponseStruct = DeleteUserResponseStruct{}

	err := gCtx.BindJSON(&req)

	if err != nil {
		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusBadRequest, res)
		return
	}

	mdlReq := Model.ModelDeleteUserRequestStruct{
		UserID: req.UserID,
	}

	err = Ctrl.Mdl.DeleteUser(mdlReq)

	if err != nil {

		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusInternalServerError, res)
		return
	}

	_, err = Ctrl.Mdl.UpdateToken(req.UserID, "")

	if err != nil {

		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusInternalServerError, res)
		return
	}

	res.IsDeleted = true
	res.IsAnyError = false

	gCtx.JSON(http.StatusOK, res)

}

func (Ctrl *ControllerStruct) EditUser(gCtx *gin.Context) {

	var req EditUserRequestStruct = EditUserRequestStruct{}

	var res EditUserResponseStruct = EditUserResponseStruct{}
	res.UserDt = req

	err := gCtx.BindJSON(&req)

	if err != nil {
		res.IsAnyError = true
		res.IsDeleted = false
		t := err.Error()
		log.Println(t)
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusBadRequest, res)
		return
	}

	mdlReq := Model.ModelEditUserRequestStruct{
		UserID:              req.UserID,
		Name:                req.Name,
		Password:            req.Password,
		Email:               req.Email,
		Is_Password_Changed: false,
	}

	err = Ctrl.Mdl.EditUser(mdlReq)

	if err != nil {

		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusInternalServerError, res)
		return
	}

	res.IsDeleted = true
	res.IsAnyError = false

	gCtx.JSON(http.StatusOK, res)

}

func (Ctrl *ControllerStruct) EditUserCred(gCtx *gin.Context) {

	var req EditUserRequestStruct = EditUserRequestStruct{}

	var res EditUserResponseStruct = EditUserResponseStruct{}
	res.UserDt = req

	err := gCtx.BindJSON(&req)

	if err != nil {
		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusBadRequest, res)
		return
	}

	mdlReq := Model.ModelEditUserRequestStruct{
		UserID:              req.UserID,
		Name:                req.Name,
		Password:            req.Password,
		Email:               req.Email,
		Is_Password_Changed: true,
	}

	err = Ctrl.Mdl.EditUser(mdlReq)

	if err != nil {

		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusInternalServerError, res)
		return
	}

	newTkn, err := Ctrl.Mdl.CreateToken(req.UserID)

	if err != nil {

		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusInternalServerError, res)
		return
	}

	_, err = Ctrl.Mdl.UpdateToken(req.UserID, newTkn)

	if err != nil {

		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusInternalServerError, res)
		return
	}

	res.IsDeleted = true
	res.IsAnyError = false

	gCtx.Header("Authorization", "Bearer "+newTkn)

	gCtx.JSON(http.StatusOK, res)

}

func (Ctrl *ControllerStruct) VerifyCred(gCtx *gin.Context) {

	var req VerifyUserRequestStruct = VerifyUserRequestStruct{}

	var res VerifyUserResponseStruct = VerifyUserResponseStruct{}

	err := gCtx.BindJSON(&req)

	if err != nil {
		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusBadRequest, res)
		return
	}

	mdlReq := Model.ModelVerifyCredRequestStruct{
		Password: req.Password,
		Email:    req.Email,
	}

	isCorrect, userID, err := Ctrl.Mdl.VerifyCred(mdlReq)

	if err != nil {

		res.IsAnyError = true
		res.IsDeleted = false
		res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
		gCtx.JSON(http.StatusInternalServerError, res)
		return
	}

	if isCorrect == true {

		newTkn, err := Ctrl.Mdl.AddToken(userID)

		if err != nil {

			res.IsAnyError = true
			res.IsDeleted = false
			res.ErrorMessages = append(res.ErrorMessages, "Internal Server Error!")
			gCtx.JSON(http.StatusInternalServerError, res)
			return
		}

		res.IsDeleted = true
		res.IsAnyError = false

		gCtx.Header("Authorization", "Bearer "+newTkn)

		gCtx.JSON(http.StatusOK, res)
	} else {

		res.IsDeleted = false
		res.IsAnyError = true

		gCtx.JSON(http.StatusNotFound, res)
	}

}

func (Ctrl *ControllerStruct) AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}

		tokenString := strings.Split(token, " ")[1]

		tkn, err := jwt.ParseWithClaims(tokenString, &Model.TokenPayLoad{}, func(token *jwt.Token) (interface{}, error) {
			return Ctrl.Mdl.Conf.JwtSecretKey, nil
		}, jwt.WithLeeway(5*time.Hour))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid or expired token"})
			return
		}

		if !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Next()
	}

}
