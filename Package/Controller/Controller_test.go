package Controller

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"
	"CarbonCreditMarketPlaceAuthAPI/Package/Configurator"
	"CarbonCreditMarketPlaceAuthAPI/Package/CustomLogger"
	"CarbonCreditMarketPlaceAuthAPI/Package/Model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type TestAddUserRequest struct {
	name     string
	email    string
	password string
}

type TestControllerStruct struct {
	suite.Suite
	Ctrl ControllerStruct
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

	recorder := httptest.NewRecorder()

	addUserRequest := TestAddUserRequest{
		name:     "Test",
		email:    "Test@test.com",
		password: "Test",
	}

	jsonData, err := json.Marshal(addUserRequest)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	reqBody := bytes.NewBuffer(jsonData)

	req1, err := http.NewRequest("POST", "/Add", reqBody)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req1)

	Ts.Assert().Equal(200, recorder.Code)

	respBody, err := ioutil.ReadAll(recorder.Body)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	Ts.Assert().Nil(err)
	Ts.Assert().NotNil(respBody)

}

func (Ts *TestControllerStruct) BeforeTest(SuiteName string, TestName string) {
	switch TestName {
	case "":
	}
}

func (Ts *TestControllerStruct) AfterTest(SuiteName string, TestName string) {
	switch TestName {
	case "":
	}
}

func (Ts *TestControllerStruct) TearDownSuite() {
	Ts.Ctrl.Conf.DB.Close()
}
