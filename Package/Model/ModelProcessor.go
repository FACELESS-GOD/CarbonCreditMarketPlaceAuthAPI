package Model

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenPayLoad struct {
	UserId    int       `json:"userid"`
	IssuedAT  time.Time `json:"issuedat"`
	ExpiredAT time.Time `json:"expiredat"`
	ID        uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func validateAddRequest(Req ModelAddUserRequestStruct) (bool, error) {

	if len(Req.Name) < 1 {
		return false, nil
	}

	if len(Req.Password) < 1 {
		return false, nil
	}

	if len(Req.Email) < 1 {
		return false, nil
	}

	var adderate rune = rune('@')
	var dot rune = rune('.')
	var isAdderPresent bool = false
	var isdotPresent bool = false
	for i := 0; i < len(Req.Email); i++ {
		if rune(Req.Email[i]) == adderate {
			isAdderPresent = true
		} else if rune(Req.Email[i]) == dot {
			isdotPresent = true
		}
	}

	if isAdderPresent != true || isdotPresent != true {
		return false, nil
	} else {
		return true, nil
	}

}

func validateEditUser(Req ModelEditUserRequestStruct) (bool, error) {
	if Req.UserID < 1 {
		return false, nil
	}

	if len(Req.Name) < 1 {
		return false, nil
	}

	var adderate rune = rune('@')
	var dot rune = rune('.')
	var isAdderPresent bool = false
	var isdotPresent bool = false
	for i := 0; i < len(Req.Email); i++ {
		if rune(Req.Email[i]) == adderate {
			isAdderPresent = true
		} else if rune(Req.Email[i]) == dot {
			isdotPresent = true
		}
	}

	if isAdderPresent != true || isdotPresent != true {
		return false, nil
	}

	if Req.Is_Password_Changed == true {
		if len(Req.Password) < 1 {
			return false, nil
		}
	}
	return true, nil
}

func validateUpdateCred(Req ModelUpdateCredRequestStruct) (bool, error) {
	if Req.UserID < 1 {
		return false, nil
	}

	if len(Req.Password) < 1 {
		return false, nil
	}

	return true, nil
}

func validateVerifyCred(Req ModelVerifyCredRequestStruct) (bool, error) {

	if len(Req.Password) < 1 {
		return false, nil
	}

	var adderate rune = rune('@')
	var dot rune = rune('.')
	var isAdderPresent bool = false
	var isdotPresent bool = false
	for i := 0; i < len(Req.Email); i++ {
		if rune(Req.Email[i]) == adderate {
			isAdderPresent = true
		} else if rune(Req.Email[i]) == dot {
			isdotPresent = true
		}
	}

	if isAdderPresent != true || isdotPresent != true {
		return false, nil
	} else {
		return true, nil
	}
}

func validateDeleteUser(Req ModelDeleteUserRequestStruct) (bool, error) {
	if Req.UserID < 1 {
		return false, nil
	}

	return true, nil
}

func (Mdl *ModelStruct) Reset() {
	Mdl.ErrorMessages = []string{}
	Mdl.IsAnyError = false
}

func (Mdl *ModelStruct) GenerateHash(Password string) (string, error) {

	var customCost int = 15
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), customCost)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(hashedPassword), nil
}

func (Mdl *ModelStruct) CreateToken(UserId int) (string, error) {

	id, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	tokenData := TokenPayLoad{
		ID:        id,
		UserId:    UserId,
		IssuedAT:  time.Now(),
		ExpiredAT: time.Now().Add(time.Duration(time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"ID":        tokenData.ID,
			"UserId":    tokenData.UserId,
			"IssuedAT":  tokenData.IssuedAT,
			"ExpiredAT": tokenData.ExpiredAT,
		})

	tokenString, err := token.SignedString([]byte(Mdl.Conf.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

const AddUserQuery string = `
INSERT INTO User (
  Name, email
) VALUES (
  ? , ? 
)
;
`

const AddUserCredQuery string = `
INSERT INTO UserCred (
  UserId, Hash_Password
) VALUES (
  ? , ? 
)
;
`

func (Mdl ModelStruct) AddUser(Req ModelAddUserRequestStruct) ModelAddUserResponseStruct {
	res := ModelAddUserResponseStruct{}

	Mdl.Reset()

	isvalid, err := validateAddRequest(Req)
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return res
	}

	if isvalid != true {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		return res
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	defer db.Commit()

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return res
	}

	response, err := db.ExecContext(ctx, AddUserQuery, Req.Name, Req.Email)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return res
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return res
		}
	}

	userID, err := response.LastInsertId()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return res
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return res
		}
	}

	password, err := Mdl.GenerateHash(Req.Password)
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return res
	}

	response, err = db.ExecContext(ctx, AddUserCredQuery, userID, password)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return res
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return res
		}
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return res
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return res
		}
	}

	log.Println(response)

	res.UserID = int(userID)

	return res
}

const DeleteUserQuery string = `
UPDATE User
SET Is_Visible = 0 , Last_Modified_Date = CURDATE()
WHERE UserId  = ? AND Is_Visible = 1
;
`

const DeleteUserCredQuery string = `
UPDATE UserCred
SET Is_Visible = 0 , Last_Modified_Date = CURDATE()
WHERE UserId  = ? AND Is_Visible = 1
;
`

func (Mdl ModelStruct) DeleteUser(Req ModelDeleteUserRequestStruct) error {

	Mdl.Reset()

	isvalid, err := validateDeleteUser(Req)

	if isvalid != true {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		return errors.New(strings.Join(Mdl.ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)
	defer db.Commit()
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return err
	}

	userCredresponse, err := db.ExecContext(ctx, DeleteUserCredQuery, Req.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return err
		}
	}

	log.Println(userCredresponse)

	userResponse, err := db.Query(DeleteUserQuery, Req.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return err
		}
	}

	log.Println(userResponse)

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return err
		}
	}

	return nil

}

const EditUserQuery string = `
UPDATE User
SET Name = ? , email = ? , Last_Modified_Date = CURDATE()
WHERE UserId  = ? AND Is_Visible = 1
;
`

const EditUserCredQuery string = `
UPDATE UserCred
SET Hash_Password = ? , Last_Modified_Date = CURDATE()
WHERE UserId  = ? AND Is_Visible = 1
;
`

func (Mdl ModelStruct) EditUser(Req ModelEditUserRequestStruct) error {
	Mdl.Reset()

	isvalid, err := validateEditUser(Req)

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return err
	}

	if isvalid != true {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		return errors.New(strings.Join(Mdl.ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)
	defer db.Commit()
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return err
	}

	response, err := db.ExecContext(ctx, EditUserQuery, Req.Name, Req.Email, Req.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return err
		}
	}

	log.Println(response)

	if Req.Is_Password_Changed == true {

		password, err := Mdl.GenerateHash(Req.Password)
		if err != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return err
		}

		response, err = db.ExecContext(ctx, EditUserCredQuery, password, Req.UserID)

		if err != nil {
			nerr := db.Rollback()
			if nerr != nil {
				Mdl.IsAnyError = true
				Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
				return nerr
			} else {
				Mdl.IsAnyError = true
				Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
				return err
			}
		}

	}

	log.Println(response)

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return err
		}
	}

	return nil
}

const UpdateUserCredQuery string = `
UPDATE UserCred
SET Hash_Password = ? , Last_Modified_Date = CURDATE()
WHERE UserId  = ? AND Is_Visible = 1
;
`

func (Mdl ModelStruct) UpdateCred(Req ModelUpdateCredRequestStruct) error {
	Mdl.Reset()

	isvalid, err := validateUpdateCred(Req)

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return err
	}

	if isvalid != true {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		return errors.New(strings.Join(Mdl.ErrorMessages, ","))
	}

	password, err := Mdl.GenerateHash(Req.Password)
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return err
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)
	defer db.Commit()
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return err
	}

	response, err := db.ExecContext(ctx, UpdateUserCredQuery, password, Req.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return err
		}
	}

	log.Println(response)

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return err
		}
	}

	return nil
}

const GetUserQuery string = `
SELECT UserId from User
WHERE email  = ? AND Is_Visible = 1
;
`

const GetUserCredQuery string = `
SELECT Hash_Password from UserCred
WHERE UserId  = ? AND Is_Visible = 1
;
`

func (Mdl ModelStruct) VerifyCred(Req ModelVerifyCredRequestStruct) (bool, int, error) {

	Mdl.Reset()

	isvalid, err := validateVerifyCred(Req)

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return false, 0, err
	}

	if isvalid != true {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		return false, 0, errors.New(strings.Join(Mdl.ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)
	defer db.Commit()
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return false, 0, err
	}

	response, err := db.Query(GetUserQuery, Req.Email)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return false, 0, nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, 0, err
		}
	}

	userId := 0

	for response.Next() {

		err = response.Scan(&userId)

		if err != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, 0, err
		}

		if userId > 1 {
			break
		}

	}

	if userId < 1 {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Wrong email.")
		response.Close()

		return false, 0, db.Rollback()
	}

	response.Close()

	response, err = db.Query(GetUserCredQuery, userId)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return false, 0, nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, 0, err
		}
	}

	var dbHashedPassword string

	for response.Next() {

		err = response.Scan(&dbHashedPassword)

		if err != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, 0, err
		}

		if len(dbHashedPassword) >= 1 {
			break
		}

	}

	if len(dbHashedPassword) < 1 {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Wrong Password.")
		response.Close()
		return false, 0, db.Rollback()
	}

	response.Close()

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return false, 0, nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbHashedPassword), []byte(Req.Password))

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return false, 0, err
	}

	return true, userId, err
}

const AddTokenQuery string = `
INSERT INTO TokenStore (
  Token, UserId
) VALUES (
  ? , ? 
)
;
`

func (Mdl ModelStruct) AddToken(UserID int) (string, error) {

	Mdl.Reset()

	if UserID < 1 {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		return "", errors.New(strings.Join(Mdl.ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	token, err := Mdl.CreateToken(UserID)

	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return "", err
	}

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)
	defer db.Commit()
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return "", err
	}

	response, err := db.Query(AddTokenQuery, token, UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return "", nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return "", err
		}
	}

	log.Println(response)

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return "", nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return "", err
		}
	}
	return token, err
}

const UpdateTokenQuery string = `
UPDATE TokenStore
SET Token = ? , Last_Modified_Date = CURDATE()
WHERE UserId  = ? AND Is_Visible = 1
;
`

func (Mdl ModelStruct) UpdateToken(UserId int, Token string) (bool, error) {

	Mdl.Reset()

	if UserId < 1 {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		return false, errors.New(strings.Join(Mdl.ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)
	defer db.Commit()
	if err != nil {
		Mdl.IsAnyError = true
		Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
		return false, err
	}

	response, err := db.Query(UpdateTokenQuery, Token, UserId)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return false, nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, err
		}
	}

	log.Println(response)

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
			return false, nerr
		} else {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, err
		}
	}
	return true, err
}

const GetTokenQuery string = `
SELECT Token from TokenStore
WHERE UserId  = ? AND Is_Visible = 1
;
`

func (Mdl ModelStruct) VerifyToken(Token string, UserID int) (bool, error) {

	Mdl.Reset()

	if UserID < 1 {
		Mdl.IsAnyError = true
		if Mdl.ErrorMessages != nil {
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		}
		return false, errors.New(strings.Join(Mdl.ErrorMessages, ","))
	}

	if len(Token) < 1 {
		Mdl.IsAnyError = true
		if Mdl.ErrorMessages != nil {
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Data is Invalid!")
		}
		return false, errors.New(strings.Join(Mdl.ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	if ctx != nil && Mdl.Conf.DB != nil {

		db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)
		defer db.Commit()
		if err != nil {
			Mdl.IsAnyError = true
			Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
			return false, err
		}

		if db != nil {

			response, err := db.Query(GetTokenQuery, UserID)

			if err != nil {
				nerr := db.Rollback()
				if nerr != nil {
					Mdl.IsAnyError = true
					Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
					return false, nerr
				} else {
					Mdl.IsAnyError = true
					Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
					return false, err
				}
			}

			dbToken := ""

			if response != nil {

				for response.Next() {
					err = response.Scan(&dbToken)

					if err != nil {
						Mdl.IsAnyError = true
						Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
						return false, err
					}

					if len(dbToken) >= 1 {
						break
					}
				}

				if len(dbToken) < 1 {
					Mdl.IsAnyError = true
					Mdl.ErrorMessages = append(Mdl.ErrorMessages, "Invalid data ")
					return false, err
				}

				err = db.Commit()

				if err != nil {
					nerr := db.Rollback()
					if nerr != nil {
						Mdl.IsAnyError = true
						Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error()+nerr.Error())
						return false, nerr
					} else {
						Mdl.IsAnyError = true
						Mdl.ErrorMessages = append(Mdl.ErrorMessages, err.Error())
						return false, err
					}
				}

				if dbToken == Token {
					return true, nil
				} else {
					Mdl.IsAnyError = true
					return false, nil
				}

			} else {
				Mdl.IsAnyError = true
				return false, nil
			}

		} else {
			Mdl.IsAnyError = true
			return false, nil
		}

	} else {
		Mdl.IsAnyError = true
		return false, nil
	}

}
