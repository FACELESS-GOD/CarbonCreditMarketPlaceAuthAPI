package Model

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func validateAddRequest(Req ModelAddUserRequestStruct) (bool, error) {

	if len(Req.Name) < 1 {
		return false, nil
	}

	if len(Req.Password) < 1 {
		return false, nil
	}

	if len(Req.email) < 1 {
		return false, nil
	}

	var adderate rune = rune('@')
	var dot rune = rune('.')

	for i := 0; i < len(Req.email); i++ {
		if rune(Req.email[i]) == adderate {
			return false, nil
		} else if rune(Req.email[i]) == dot {
			return false, nil
		}
	}

	return true, nil
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

	for i := 0; i < len(Req.email); i++ {
		if rune(Req.email[i]) == adderate {
			return false, nil
		} else if rune(Req.email[i]) == dot {
			return false, nil
		}
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

	for i := 0; i < len(Req.email); i++ {
		if rune(Req.email[i]) == adderate {
			return false, nil
		} else if rune(Req.email[i]) == dot {
			return false, nil
		}
	}

	return true, nil
}

func validateDeleteUser(Req ModelDeleteUserRequestStruct) (bool, error) {
	if Req.UserID < 1 {
		return false, nil
	}

	return true, nil
}

func (Mdl *ModelStruct) Reset() {
	ErrorMessages = []string{}
	IsAnyError = false
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

func (Mdl *ModelStruct) createToken(UserId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"UserId": UserId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(Mdl.Conf.JwtSecretKey)
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

func (Mdl *ModelStruct) AddUser(Req ModelAddUserRequestStruct) ModelAddUserResponseStruct {
	res := ModelAddUserResponseStruct{}

	Mdl.Reset()

	isvalid, err := validateAddRequest(Req)
	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return res
	}

	if isvalid != true {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return res
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return res
	}

	response, err := db.ExecContext(ctx, AddUserQuery, Req.Name, Req.email)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return res
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return res
		}
	}

	userID, err := response.LastInsertId()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return res
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return res
		}
	}

	response, err = db.ExecContext(ctx, AddUserCredQuery, userID, Req.Password)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return res
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return res
		}
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return res
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return res
		}
	}

	log.Println(response)

	res.UserID = int(userID)

	return res
}

const DeleteUserQuery string = `
UPDATE User
SET Is_Visible = 0 , Last_Modified_Date = GETDATE()
WHERE UserId  = ?
;
`

const DeleteUserCredQuery string = `
UPDATE UserCred
SET Is_Visible = 0 , Last_Modified_Date = GETDATE()
WHERE UserId  = ?
;
`

func (Mdl *ModelStruct) DeleteUser(Req ModelDeleteUserRequestStruct) error {

	Mdl.Reset()

	isvalid, err := validateDeleteUser(Req)

	if isvalid != true {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return errors.New(strings.Join(ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	response, err := db.ExecContext(ctx, DeleteUserQuery, Req.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return err
		}
	}

	log.Println(response)

	response, err = db.ExecContext(ctx, DeleteUserCredQuery, Req.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return err
		}
	}

	log.Println(response)

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return err
		}
	}

	return nil

}

const EditUserQuery string = `
UPDATE User
SET Name = ? , email = ? , Last_Modified_Date = GETDATE()
WHERE UserId  = ?
;
`

const EditUserCredQuery string = `
UPDATE UserCred
SET Hash_Password = ? , Last_Modified_Date = GETDATE()
WHERE UserId  = ?
;
`

func (Mdl *ModelStruct) EditUser(Req ModelEditUserRequestStruct) error {
	Mdl.Reset()

	isvalid, err := validateEditUser(Req)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	if isvalid != true {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return errors.New(strings.Join(ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	response, err := db.ExecContext(ctx, EditUserQuery, Req.Name, Req.email, Req.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return err
		}
	}

	log.Println(response)

	if Req.Is_Password_Changed == true {

		password, err := Mdl.GenerateHash(Req.Password)
		if err != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return err
		}

		response, err = db.ExecContext(ctx, EditUserCredQuery, password, Req.UserID)

		if err != nil {
			nerr := db.Rollback()
			if nerr != nil {
				IsAnyError = true
				ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
				return nerr
			} else {
				IsAnyError = true
				ErrorMessages = append(ErrorMessages, err.Error())
				return err
			}
		}

	}

	log.Println(response)

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return err
		}
	}

	return nil
}

const UpdateUserCredQuery string = `
UPDATE UserCred
SET Hash_Password = ? , Last_Modified_Date = GETDATE()
WHERE UserId  = ?
;
`

func (Mdl *ModelStruct) UpdateCred(Req ModelUpdateCredRequestStruct) error {
	Mdl.Reset()

	isvalid, err := validateUpdateCred(Req)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	if isvalid != true {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return errors.New(strings.Join(ErrorMessages, ","))
	}

	password, err := Mdl.GenerateHash(Req.Password)
	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	response, err := db.ExecContext(ctx, UpdateUserCredQuery, password, Req.UserID)

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return err
		}
	}

	log.Println(response)

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return err
		}
	}

	return nil
}

const GetUserQuery string = `
SELECT UserId from User
WHERE email  = ?
;
`

const GetUserCredQuery string = `
SELECT Hash_Password from UserCred
WHERE UserId  = ?
;
`

func (Mdl *ModelStruct) VerifyCred(Req ModelVerifyCredRequestStruct) (bool, error) {

	Mdl.Reset()

	isvalid, err := validateVerifyCred(Req)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	if isvalid != true {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return false, errors.New(strings.Join(ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	response := db.QueryRowContext(ctx, GetUserQuery, Req.email)

	if response.Err() != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error())
			return false, err
		}
	}

	var userId int

	err = response.Scan(&userId)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	response = db.QueryRowContext(ctx, GetUserCredQuery, userId)

	if response.Err() != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error())
			return false, err
		}
	}

	var dbHashedPassword string

	err = response.Scan(&dbHashedPassword)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return false, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbHashedPassword), []byte(Req.Password))

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	return true, err
}

const AddTokenQuery string = `
INSERT INTO UserCred (
  Token, UserId
) VALUES (
  ? , ? 
)
;
`

func (Mdl *ModelStruct) AddToken(UserID int) (bool, error) {

	Mdl.Reset()

	if UserID < 1 {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return false, errors.New(strings.Join(ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	token, err := Mdl.createToken(UserID)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	response := db.QueryRowContext(ctx, AddTokenQuery, token, UserID)

	if response.Err() != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error())
			return false, err
		}
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return false, err
		}
	}
	return true, err
}

const UpdateTokenQuery string = `
UPDATE TokenStore
SET Token = ? , Last_Modified_Date = GETDATE()
WHERE UserId  = ?
;
`

func (Mdl *ModelStruct) UpdateToken(UserId int, Token string) (bool, error) {

	Mdl.Reset()

	if UserId < 1 {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return false, errors.New(strings.Join(ErrorMessages, ","))
	}

	if len(Token) < 1 {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return false, errors.New(strings.Join(ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	response := db.QueryRowContext(ctx, UpdateTokenQuery, Token, UserId)

	if response.Err() != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error())
			return false, err
		}
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return false, err
		}
	}
	return true, err
}

const GetTokenQuery string = `
SELECT Hash_Password from UserCred
WHERE UserId  = ?
;
`

func (Mdl *ModelStruct) VerifyToken(Token string, UserID int) (bool, error) {

	Mdl.Reset()

	if UserID < 1 {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return false, errors.New(strings.Join(ErrorMessages, ","))
	}

	if len(Token) < 1 {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, "Data is Invalid!")
		return false, errors.New(strings.Join(ErrorMessages, ","))
	}

	ctx := context.WithoutCancel(context.Background())

	db, err := Mdl.Conf.DB.BeginTx(ctx, &Mdl.Conf.TxOption)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	response := db.QueryRowContext(ctx, GetTokenQuery, UserID)

	if response.Err() != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, response.Err().Error())
			return false, err
		}
	}

	var dbToken string

	err = response.Scan(&dbToken)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return false, err
	}

	err = db.Commit()

	if err != nil {
		nerr := db.Rollback()
		if nerr != nil {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error()+nerr.Error())
			return false, nerr
		} else {
			IsAnyError = true
			ErrorMessages = append(ErrorMessages, err.Error())
			return false, err
		}
	}

	if dbToken == Token {
		return true, nil
	} else {
		return false, nil
	}
}
