package Model

import (
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
)

type ModelInterface interface {
	AddUser(ModelAddUserRequestStruct) ModelAddUserResponseStruct
	DeleteUser(ModelDeleteUserRequestStruct)
	EditUser(ModelEditUserRequestStruct)
	UpdateCred(ModelUpdateCredRequestStruct)
	VerifyToken(Token string, UserID int) (bool, error)
	AddToken(UserID int) (bool, error)
	UpdateToken(UserID int, Token string) (bool, error)
	VerifyCred(ModelVerifyCredRequestStruct) (bool, error)
}

var ErrorMessages []string
var IsAnyError bool

type ModelStruct struct {
	Conf Configurator.ConfiguratorStruct
	Log  CustomLogger.CustomLoggerInterface
}

func NewModel(Conf Configurator.ConfiguratorStruct, Log CustomLogger.CustomLoggerInterface) ModelStruct {
	mdl := ModelStruct{}
	mdl.Conf = Conf
	mdl.Log = Log
	return mdl
}

type ModelAddUserRequestStruct struct {
	Name     string
	email    string
	Password string
}
type ModelAddUserResponseStruct struct {
	UserID int
}
type ModelDeleteUserRequestStruct struct {
	UserID int
}
type ModelEditUserRequestStruct struct {
	UserID              int
	Name                string
	email               string
	Is_Password_Changed bool
	Password            string
}
type ModelUpdateCredRequestStruct struct {
	Password string
	UserID   int
}

type ModelVerifyTokenRequestStruct struct {
	Token  string
	UserID int
}

type ModelVerifyCredRequestStruct struct {
	email    string
	Password string
}
