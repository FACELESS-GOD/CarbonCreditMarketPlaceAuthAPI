package Model

import (
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
)

type ModelInterface interface {
	AddUser(ModelAddUserRequestStruct) ModelAddUserResponseStruct
	DeleteUser(ModelDeleteUserRequestStruct) error
	EditUser(ModelEditUserRequestStruct) error
	UpdateCred(ModelUpdateCredRequestStruct) error
	VerifyToken(Token string, UserID int) (bool, error)
	AddToken(UserID int) (string, error)
	UpdateToken(UserID int, Token string) (bool, error)
	VerifyCred(ModelVerifyCredRequestStruct) (bool, int, error)
}

type ModelStruct struct {
	Conf          Configurator.ConfiguratorStruct
	Log           CustomLogger.CustomLoggerInterface
	ErrorMessages []string
	IsAnyError    bool
}

func NewModel(Conf Configurator.ConfiguratorStruct, Log CustomLogger.CustomLoggerInterface) ModelStruct {
	mdl := ModelStruct{}
	mdl.Conf = Conf
	mdl.Log = Log
	mdl.ErrorMessages = []string{}
	return mdl
}

type ModelAddUserRequestStruct struct {
	Name     string
	Email    string
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
	Email               string
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
	Email    string
	Password string
}
