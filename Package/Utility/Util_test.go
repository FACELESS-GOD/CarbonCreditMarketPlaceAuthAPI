package Util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestUtil struct {
	suite.Suite
	ut Util
}

func TestUtils(t *testing.T) {
	suite.Run(t, &TestUtil{})
}

func (Ts *TestUtil) TestRandomNumber() {
	num, err := Ts.ut.RandomNumber(0)
	Ts.Require().NotNil(err)
	num, err = Ts.ut.RandomNumber(10000000000000000)
	Ts.Require().Nil(err)
	fmt.Println(num)
}

func (Ts *TestUtil) TestRandomString() {
	st, err := Ts.ut.RandomString(0)
	Ts.Require().NotNil(err)
	st, err = Ts.ut.RandomString(10000000000000000)
	Ts.Require().NotNil(err)
	st, err = Ts.ut.RandomString(10)
	Ts.Require().Nil(err)
	Ts.Require().Equal(len(st) >= 1, true)
	fmt.Println(st)
}

func (Ts *TestUtil) SetupSuite() {
	Ts.ut = Util{}
}

func (Ts *TestUtil) TearDownSuite() {
	//Ts.ut = nil
}
