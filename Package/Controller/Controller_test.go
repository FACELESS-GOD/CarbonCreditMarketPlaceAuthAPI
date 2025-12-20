package Controller

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type TestAddUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TestAddUserResponseStruct struct {
	UserID int `json:"userid"`
}

type TestControllerStruct struct {
	suite.Suite
	Ctrl         ControllerStruct
	DeleteUserID int
	EditUserID   int
	EditCredUsr  Model.ModelAddUserRequestStruct
}

func (Ctrl *TestControllerStruct) reset() {
	Ctrl.DeleteUserID = 0
	Ctrl.EditUserID = 0
	Ctrl.EditCredUsr = Model.ModelAddUserRequestStruct{}
}

func TestTestControllerstruct(testObj *testing.T) {
	suite.Run(testObj, &TestControllerStruct{})
}

func (Ts *TestControllerStruct) Reset() {

}

func (Ts *TestControllerStruct) SetupSuite() {
	conf, err := Configurator.NewConfiguration(DevMode.Client)

	if err != nil {
		Ts.FailNow("DBConn Failed", err)
	}

	logger, err := CustomLogger.NewLogger(DevMode.Client)

	if err != nil {
		Ts.FailNow("Logger Failed", err)
	}

	mdl := Model.NewModel(conf, logger)

	ctrl := NewCtrl(mdl, conf, nil)

	Ts.Ctrl = ctrl

}

func (Ts *TestControllerStruct) TestAddUser() {

	router := gin.Default()

	router.POST("/Add", Ts.Ctrl.AddUser)

	recorder := httptest.NewRecorder()

	addUserRequest := TestAddUserRequest{
		Name:     "Test",
		Email:    "Test@test.com",
		Password: "Test",
	}

	jsonData, err := json.Marshal(&addUserRequest)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	//reqBody := bytes.NewBuffer(jsonData)
	req1, _ := http.NewRequest("POST", "/Add", strings.NewReader(string(jsonData)))

	///req1, err := http.NewRequest("POST", "/Add", reqBody)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req1)

	Ts.Assert().Equal(200, recorder.Code)

	bodyBytes, err := io.ReadAll(recorder.Body)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	respBody := TestAddUserResponseStruct{}

	err = json.Unmarshal(bodyBytes, &respBody)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	hdr := recorder.Header()

	Ts.Assert().Nil(err)
	Ts.Assert().NotNil(respBody)
	Ts.Assert().NotNil(hdr["Authorization"])
	Ts.Assert().NotNil(respBody.UserID)

}

func (Ts *TestControllerStruct) TestDeleteUser() {

	router := gin.Default()

	router.DELETE("/Delete", Ts.Ctrl.DeleteUser)

	recorder := httptest.NewRecorder()

	deleteUserRequestStruct := DeleteUserRequestStruct{
		UserID: Ts.DeleteUserID,
	}

	jsonData, err := json.Marshal(&deleteUserRequestStruct)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	req1, _ := http.NewRequest("DELETE", "/Delete", strings.NewReader(string(jsonData)))

	if err != nil {
		Ts.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req1)

	Ts.Assert().Equal(200, recorder.Code)

	bodyBytes, err := io.ReadAll(recorder.Body)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	respBody := DeleteUserResponseStruct{}

	err = json.Unmarshal(bodyBytes, &respBody)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	Ts.Assert().Nil(err)
	Ts.Assert().NotNil(respBody)
	Ts.Assert().Equal(respBody.IsAnyError, false)
	Ts.Assert().Equal(respBody.IsDeleted, true)

	Ts.Assert().Nil(respBody.ErrorMessages)
	Ts.Assert().Equal(len(respBody.ErrorMessages) < 1, true)

}

func (Ts *TestControllerStruct) TestEditUser() {

	router := gin.Default()

	router.PUT("/EditUsr", Ts.Ctrl.EditUser)

	recorder := httptest.NewRecorder()

	EditUserRequestStruct := EditUserRequestStruct{
		UserID:            Ts.EditUserID,
		Name:              "Test2",
		Email:             "Test2@TEst2.com",
		Ispasswordchanged: false,
		Password:          "NewPassword",
	}

	jsonData, err := json.Marshal(&EditUserRequestStruct)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	req1, _ := http.NewRequest("PUT", "/EditUsr", strings.NewReader(string(jsonData)))

	if err != nil {
		Ts.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req1)

	Ts.Assert().Equal(200, recorder.Code)

	bodyBytes, err := io.ReadAll(recorder.Body)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	respBody := EditUserResponseStruct{}

	err = json.Unmarshal(bodyBytes, &respBody)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	Ts.Assert().Nil(err)
	Ts.Assert().NotNil(respBody)
	Ts.Assert().Equal(respBody.IsAnyError, false)
	Ts.Assert().Equal(respBody.IsDeleted, true)

	Ts.Assert().Nil(respBody.ErrorMessages)
	Ts.Assert().Equal(len(respBody.ErrorMessages) < 1, true)

}

func (Ts *TestControllerStruct) TestEditUserCred() {

	router := gin.Default()

	router.PUT("/EditUsrPass", Ts.Ctrl.EditUserCred)

	recorder := httptest.NewRecorder()

	EditUserRequestStruct := EditUserRequestStruct{
		UserID:            Ts.EditUserID,
		Name:              Ts.EditCredUsr.Name,
		Email:             Ts.EditCredUsr.Email,
		Ispasswordchanged: true,
		Password:          "NewPassword",
	}

	jsonData, err := json.Marshal(&EditUserRequestStruct)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	req1, _ := http.NewRequest("PUT", "/EditUsrPass", strings.NewReader(string(jsonData)))

	if err != nil {
		Ts.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req1)

	Ts.Assert().Equal(200, recorder.Code)

	bodyBytes, err := io.ReadAll(recorder.Body)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	respBody := EditUserResponseStruct{}

	err = json.Unmarshal(bodyBytes, &respBody)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	hdr := recorder.Header()

	Ts.Assert().Nil(err)
	Ts.Assert().NotNil(respBody)

	Ts.Assert().NotNil(hdr["Authorization"])
	Ts.Assert().Equal(respBody.IsAnyError, false)
	Ts.Assert().Equal(respBody.IsDeleted, true)

	Ts.Assert().Nil(respBody.ErrorMessages)
	Ts.Assert().Equal(len(respBody.ErrorMessages) < 1, true)

}

func (Ts *TestControllerStruct) TestVerifyUser() {

	router := gin.Default()

	router.GET("/Verify", Ts.Ctrl.VerifyCred)

	recorder := httptest.NewRecorder()

	VerifyUserRequestStruct := VerifyUserRequestStruct{
		Email:    Ts.EditCredUsr.Email,
		Password: Ts.EditCredUsr.Password,
	}

	jsonData, err := json.Marshal(&VerifyUserRequestStruct)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	req1, _ := http.NewRequest("GET", "/Verify", strings.NewReader(string(jsonData)))

	if err != nil {
		Ts.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req1)

	Ts.Assert().Equal(200, recorder.Code)

	bodyBytes, err := io.ReadAll(recorder.Body)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	respBody := VerifyUserResponseStruct{}

	err = json.Unmarshal(bodyBytes, &respBody)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	hdr := recorder.Header()

	Ts.Assert().Nil(err)
	Ts.Assert().NotNil(respBody)

	Ts.Assert().NotNil(hdr["Authorization"])
	Ts.Assert().Equal(respBody.IsAnyError, false)
	Ts.Assert().Equal(respBody.IsDeleted, true)

	Ts.Assert().Nil(respBody.ErrorMessages)
	Ts.Assert().Equal(len(respBody.ErrorMessages) < 1, true)

}

func (Ts *TestControllerStruct) BeforeTest(SuiteName string, TestName string) {
	switch TestName {
	case "TestDeleteUser":
		Ts.reset()
		usrDT := Model.ModelAddUserRequestStruct{
			Name:     "Test",
			Email:    "Test@Test.com",
			Password: "Test",
		}
		res := Ts.Ctrl.Mdl.AddUser(usrDT)

		if Ts.Ctrl.Mdl.IsAnyError == true {
			Ts.FailNow("Can't add a new User for testing TestDeleteUser")
		}
		Ts.DeleteUserID = res.UserID

	case "TestEditUser":
		Ts.reset()
		usrDT := Model.ModelAddUserRequestStruct{
			Name:     "Test",
			Email:    "Test@Test.com",
			Password: "Test",
		}
		res := Ts.Ctrl.Mdl.AddUser(usrDT)

		if Ts.Ctrl.Mdl.IsAnyError == true {
			Ts.FailNow("Can't add a new User for testing TestDeleteUser")
		}
		Ts.EditUserID = res.UserID

	case "TestEditUserCred":
		Ts.reset()
		usrDT := Model.ModelAddUserRequestStruct{
			Name:     "Test",
			Email:    "Test@Test.com",
			Password: "Test",
		}
		res := Ts.Ctrl.Mdl.AddUser(usrDT)

		if Ts.Ctrl.Mdl.IsAnyError == true {
			Ts.FailNow("Can't add a new User for testing TestDeleteUser")
		}
		Ts.EditCredUsr = usrDT
		Ts.EditUserID = res.UserID

	case "TestVerifyUser":
		Ts.reset()
		usrDT := Model.ModelAddUserRequestStruct{
			Name:     "TestMohit",
			Email:    "TestMohit@Test.com",
			Password: "Test",
		}
		res := Ts.Ctrl.Mdl.AddUser(usrDT)

		if Ts.Ctrl.Mdl.IsAnyError == true {
			Ts.FailNow("Can't add a new User for testing TestDeleteUser")
		}
		Ts.EditCredUsr = usrDT
		Ts.EditUserID = res.UserID
	}
}

func (Ts *TestControllerStruct) AfterTest(SuiteName string, TestName string) {
	switch TestName {
	case "TestDeleteUser":
		Ts.Ctrl.Mdl.DeleteUser(Model.ModelDeleteUserRequestStruct{UserID: Ts.DeleteUserID})
	case "TestEditUser":
		Ts.Ctrl.Mdl.DeleteUser(Model.ModelDeleteUserRequestStruct{UserID: Ts.EditUserID})
	case "TestEditUserCred":
		Ts.Ctrl.Mdl.DeleteUser(Model.ModelDeleteUserRequestStruct{UserID: Ts.EditUserID})

	case "TestVerifyUser":
		Ts.Ctrl.Mdl.DeleteUser(Model.ModelDeleteUserRequestStruct{UserID: Ts.EditUserID})
	}
}

func (Ts *TestControllerStruct) TearDownSuite() {
	Ts.Ctrl.Conf.DB.Close()
}
