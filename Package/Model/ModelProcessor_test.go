package Model

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TestModelstruct struct {
	suite.Suite
	Mdl                 ModelStruct
	AddUserID           int
	ToBeDeletedUserID   int
	EditUser            ModelAddUserRequestStruct
	EditUserID          int
	UpdateCredID        int
	TokenUserID         int
	CredUser            ModelAddUserRequestStruct
	CredUserID          int
	VerifyTokenUser     ModelAddUserRequestStruct
	VerifyTokenUserID   int
	VerifyOriginalToken string
	UpdateTokenUser     ModelAddUserRequestStruct
	UpdateTokenUserID   int
}

func (Ist *TestModelstruct) Reset() {
	Ist.AddUserID = 0
	Ist.ToBeDeletedUserID = 0
	Ist.EditUserID = 0
	Ist.UpdateCredID = 0
	Ist.TokenUserID = 0
	Ist.CredUserID = 0
	Ist.VerifyTokenUserID = 0
	Ist.UpdateTokenUserID = 0
	Ist.VerifyOriginalToken = ""
	Ist.EditUser = ModelAddUserRequestStruct{}
	Ist.CredUser = ModelAddUserRequestStruct{}
	Ist.VerifyTokenUser = ModelAddUserRequestStruct{}
	Ist.UpdateTokenUser = ModelAddUserRequestStruct{}
}

func TestTestModelstruct(testObj *testing.T) {
	suite.Run(testObj, &TestModelstruct{})
}

func (Its *TestModelstruct) SetupSuite() {
	conf, err := Configurator.NewConfiguration(DevMode.Test, "../../")

	if err != nil {
		Its.FailNow("DBConn Failed", err)
	}

	logger, err := CustomLogger.NewLogger(DevMode.Test)

	if err != nil {
		Its.FailNow("DBConn Failed", err)
	}

	mdl := NewModel(conf, logger)

	Its.Mdl = mdl
}

func (Its *TestModelstruct) TestAddUser() {
	Its.Reset()

	req := ModelAddUserRequestStruct{}
	req.Email = ""
	req.Name = ""
	req.Password = ""

	Its.Mdl.AddUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "q@q.co"
	req.Name = ""
	req.Password = ""

	Its.Mdl.AddUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "q@q.co"
	req.Name = "q"
	req.Password = ""

	Its.Mdl.AddUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "qqco"
	req.Name = "q"
	req.Password = "q"

	Its.Mdl.AddUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "q@qco"
	req.Name = "q"
	req.Password = "q"

	Its.Mdl.AddUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "qq.co"
	req.Name = "q"
	req.Password = "q"

	Its.Mdl.AddUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "q@q.co"
	req.Name = "q"
	req.Password = "q"

	res := Its.Mdl.AddUser(req)

	Its.AddUserID = res.UserID

	Its.Require().Falsef(Its.Mdl.IsAnyError, "Test Case Passes!")

}

func (Its *TestModelstruct) TestEditUser() {

	req := ModelEditUserRequestStruct{}
	req.Email = ""
	req.Name = ""
	req.Password = ""

	Its.Mdl.EditUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req = ModelEditUserRequestStruct{}
	req.UserID = Its.EditUserID
	req.Email = ""
	req.Name = ""
	req.Password = ""

	Its.Mdl.EditUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.UserID = Its.EditUserID
	req.Email = "q@q.co"
	req.Name = ""
	req.Password = ""

	Its.Mdl.EditUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.UserID = Its.EditUserID
	req.Email = "q@q.co"
	req.Name = ""
	req.Password = ""
	req.Is_Password_Changed = true

	Its.Mdl.EditUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.UserID = Its.EditUserID
	req.Email = "q@q.co"
	req.Name = "q"
	req.Password = ""
	req.Is_Password_Changed = true

	Its.Mdl.EditUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.UserID = Its.EditUserID
	req.Email = "qqco"
	req.Name = "q"
	req.Password = "q"
	req.Is_Password_Changed = true

	Its.Mdl.EditUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.UserID = Its.EditUserID
	req.Email = "q@qco"
	req.Name = "q"
	req.Password = "q"
	req.Is_Password_Changed = true

	Its.Mdl.EditUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.UserID = Its.EditUserID
	req.Email = "qq.co"
	req.Name = "q"
	req.Password = "q"
	req.Is_Password_Changed = true

	Its.Mdl.EditUser(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Is_Password_Changed = true
	req.UserID = Its.EditUserID
	req.Email = "q@q.co"
	req.Name = "q"
	req.Password = "q"

	Its.Mdl.EditUser(req)

	//Its.AddUserID = res.UserID

	Its.Require().Falsef(Its.Mdl.IsAnyError, "Test Case Passes!")

}

func (Its *TestModelstruct) TestDeleteUser() {

	req := ModelDeleteUserRequestStruct{}
	req.UserID = 0

	Its.Mdl.DeleteUser(req)

	Its.Require().True(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.UserID = Its.ToBeDeletedUserID

	Its.Mdl.DeleteUser(req)

	Its.Require().False(Its.Mdl.IsAnyError)

	justInitializedReq := ModelDeleteUserRequestStruct{}

	Its.Mdl.DeleteUser(justInitializedReq)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])
}

func (Its *TestModelstruct) TestUpdateCred() {

	req := ModelUpdateCredRequestStruct{}
	req.UserID = 0
	req.Password = ""

	Its.Mdl.UpdateCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req = ModelUpdateCredRequestStruct{}
	req.UserID = Its.UpdateCredID
	req.Password = ""

	Its.Mdl.UpdateCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.UserID = Its.UpdateCredID
	req.Password = "Test"

	Its.Mdl.UpdateCred(req)

	Its.Require().Falsef(Its.Mdl.IsAnyError, "Test Case Passes!")

}

func (Its *TestModelstruct) TestVerifyCred() {

	req := ModelVerifyCredRequestStruct{}
	req.Email = ""
	req.Password = ""

	Its.Mdl.VerifyCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "q"
	req.Password = ""

	Its.Mdl.VerifyCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "q@"
	req.Password = ""

	Its.Mdl.VerifyCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "q."
	req.Password = ""

	Its.Mdl.VerifyCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = Its.CredUser.Email
	req.Password = ""

	Its.Mdl.VerifyCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = "pdoaiuf@apsdofb.ikaUJSDGFD"
	req.Password = Its.CredUser.Password

	Its.Mdl.VerifyCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = Its.CredUser.Email
	req.Password = "lakjshjkcbv"

	Its.Mdl.VerifyCred(req)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	req.Email = Its.CredUser.Email
	req.Password = Its.CredUser.Password

	Its.Mdl.VerifyCred(req)

	Its.Require().Falsef(Its.Mdl.IsAnyError, "Test Case Passes!")

}

func (Its *TestModelstruct) TestAddToken() {

	Its.Mdl.AddToken(0)

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	token, err := Its.Mdl.AddToken(Its.TokenUserID)

	Its.Require().Falsef(Its.Mdl.IsAnyError, "Test Case Passes!")

	Its.Require().NotNil(token)
	Its.Require().Nil(err, "Test Case Passes!")

}

func (Its *TestModelstruct) TestVerifyToken() {

	Its.Mdl.VerifyToken("", 0)

	Its.Require().True(Its.Mdl.IsAnyError)

	Its.Mdl.VerifyToken("", 10)

	Its.Require().True(Its.Mdl.IsAnyError)

	Its.Mdl.VerifyToken("test", 0)

	Its.Require().True(Its.Mdl.IsAnyError)

	Its.Mdl.VerifyToken("987654321206549819", Its.VerifyTokenUserID)

	Its.Require().True(Its.Mdl.IsAnyError)

	Its.Mdl.VerifyToken(Its.VerifyOriginalToken, Its.VerifyTokenUserID)

	Its.Require().False(Its.Mdl.IsAnyError)

}

func (Its *TestModelstruct) TestUpdateToken() {

	Its.Mdl.UpdateToken(0, "")

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	Its.Mdl.UpdateToken(0, "test")

	Its.Require().Truef(Its.Mdl.IsAnyError, Its.Mdl.ErrorMessages[0])

	tkn, err := Its.Mdl.CreateToken(Its.UpdateTokenUserID, time.Now().Add(time.Duration(time.Hour)))

	if err != nil {
		Its.FailNow(err.Error())
	}

	Its.Mdl.UpdateToken(Its.UpdateTokenUserID, tkn)

	Its.Require().Falsef(Its.Mdl.IsAnyError, "Test Case Passes!")

}

func (Its *TestModelstruct) BeforeTest(SuiteName string, TestName string) {
	switch TestName {
	case "TestDeleteUser":
		Its.Reset()
		req := ModelAddUserRequestStruct{}
		req.Email = "Test@Test.com"
		req.Name = "TestName"
		req.Password = "TestPassword"

		Res := Its.Mdl.AddUser(req)
		Its.ToBeDeletedUserID = Res.UserID

	case "TestEditUser":
		Its.Reset()
		req := ModelAddUserRequestStruct{}
		req.Email = "Test@Test.com"
		req.Name = "TestName"
		req.Password = "TestPassword"

		Res := Its.Mdl.AddUser(req)
		Its.EditUser = req
		Its.EditUserID = Res.UserID

	case "TestUpdateCred":
		Its.Reset()
		req := ModelAddUserRequestStruct{}
		req.Email = "Test@Test.com"
		req.Name = "TestName"
		req.Password = "TestPassword"

		Res := Its.Mdl.AddUser(req)
		Its.UpdateCredID = Res.UserID

	case "TestAddToken":

		Its.Reset()
		req := ModelAddUserRequestStruct{}
		req.Email = "Test@Test.com"
		req.Name = "TestName"
		req.Password = "TestPassword"

		Res := Its.Mdl.AddUser(req)
		Its.TokenUserID = Res.UserID

	case "TestVerifyCred":

		Its.Reset()
		req := ModelAddUserRequestStruct{}
		req.Email = "Test@Test.com"
		req.Name = "TestName"
		req.Password = "TestPassword"

		Res := Its.Mdl.AddUser(req)

		Its.CredUserID = Res.UserID
		Its.CredUser = req

	case "TestVerifyToken":

		Its.Reset()
		req := ModelAddUserRequestStruct{}
		req.Email = "Test@Test.com"
		req.Name = "TestName"
		req.Password = "TestPassword"

		Res := Its.Mdl.AddUser(req)

		Its.VerifyTokenUser = req
		Its.VerifyTokenUserID = Res.UserID

		tkn, err := Its.Mdl.AddToken(Its.VerifyTokenUserID)
		if err != nil {
			Its.FailNow(err.Error())
		}
		Its.VerifyOriginalToken = tkn

	case "TestUpdateToken":

		Its.Reset()
		req := ModelAddUserRequestStruct{}
		req.Email = "Test@Test.com"
		req.Name = "TestName"
		req.Password = "TestPassword"

		Res := Its.Mdl.AddUser(req)

		Its.UpdateTokenUser = req
		Its.UpdateTokenUserID = Res.UserID

		tkn, err := Its.Mdl.AddToken(Its.UpdateTokenUserID)
		if err != nil {
			Its.FailNow(err.Error())
		}
		Its.VerifyOriginalToken = tkn
	}
}

func (Its *TestModelstruct) AfterTest(SuiteName string, TestName string) {

	switch TestName {
	case "TestAddUser":
		req := ModelDeleteUserRequestStruct{UserID: Its.AddUserID}
		Its.Mdl.DeleteUser(req)

	case "TestEditUser":
		req := ModelDeleteUserRequestStruct{UserID: Its.EditUserID}
		Its.Mdl.DeleteUser(req)

	case "TestUpdateCred":
		req := ModelDeleteUserRequestStruct{UserID: Its.UpdateCredID}
		Its.Mdl.DeleteUser(req)

	case "TestAddToken":
		req := ModelDeleteUserRequestStruct{UserID: Its.TokenUserID}
		Its.Mdl.DeleteUser(req)

	case "TestVerifyCred":
		req := ModelDeleteUserRequestStruct{UserID: Its.CredUserID}
		Its.Mdl.DeleteUser(req)

	case "TestVerifyToken":
		req := ModelDeleteUserRequestStruct{UserID: Its.VerifyTokenUserID}
		Its.Mdl.DeleteUser(req)

	case "TestUpdateToken":
		req := ModelDeleteUserRequestStruct{UserID: Its.UpdateTokenUserID}
		Its.Mdl.DeleteUser(req)

	}
}

func (Its *TestModelstruct) TearDownSuite() {
	Its.Mdl.Conf.DB.Close()
}
