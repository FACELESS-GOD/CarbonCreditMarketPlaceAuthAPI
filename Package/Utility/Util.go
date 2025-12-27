package Util

import (
	"errors"
	"log"
	"math/rand/v2"
)

type Util struct{}

func (Ut *Util) RandomNumber(Salt int) (int32, error) {
	if Salt < 1 {
		return 0, errors.New("RandomSalt is less than 1")
	}
	return rand.Int32N(int32(Salt)), nil
}

func (Ut *Util) RandomString(num int) (string, error) {

	if num < 1 {
		return "", errors.New("Length of string cannot be less than 1")
	}

	if num > 254 {
		return "", errors.New("Length of string cannot be greated than 256")
	}

	base := "QWERTYUIOPASDFGHJKLZXCVBNM1234567890"

	baseRandomString := []byte{}

	for i := 1; i <= num; i++ {
		newNumIndex, err := Ut.RandomNumber(len(base))
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		baseRandomString = append(baseRandomString, base[newNumIndex])
	}

	return string(baseRandomString), nil

}
