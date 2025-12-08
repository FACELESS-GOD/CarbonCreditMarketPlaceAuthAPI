package Model

import (
	"context"
	"errors"
	"log"
	"strings"
)

/*

	AddUser(ModelAddUserRequestStruct) (ModelAddUserResponseStruct, error)
	DeleteUser(ModelDeleteUserRequestStruct) error
	EditUser(ModelEditUserRequestStruct) error
	UpdateCred(ModelUpdateCredRequestStruct) error
	VerifyToken(ModelVerifyTokenRequestStruct) (bool, error)
	VerifyCred(ModelVerifyCredRequestStruct) (bool, error)

*/

func (Mdl *ModelStruct) ValidationResponse(Mode int) (bool, error) {
	switch Mode {
	case add:
		return true, nil
	case delete:
		return true, nil
	case edit:
		return true, nil
	case update:
		return true, nil
	case verifyToken:
		return true, nil
	case verifyCred:
		return true, nil
	}
	return true, nil

}

func (Mdl *ModelStruct) Reset() {
	ErrorMessages = []string{}
	IsAnyError = false
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

	isValid, err := Mdl.ValidationResponse(add)
	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return res
	}

	if isValid != true {
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

	isValid, err := Mdl.ValidationResponse(delete)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	if isValid != true {
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

	isValid, err := Mdl.ValidationResponse(edit)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	if isValid != true {
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

	response, err = db.ExecContext(ctx, EditUserCredQuery, Req.Password, Req.UserID)

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

const UpdateUserCredQuery string = `
UPDATE UserCred
SET Hash_Password = ? , Last_Modified_Date = GETDATE()
WHERE UserId  = ?
;
`

func (Mdl *ModelStruct) UpdateCred(Req ModelUpdateCredRequestStruct) error {
	Mdl.Reset()

	isValid, err := Mdl.ValidationResponse(editCred)

	if err != nil {
		IsAnyError = true
		ErrorMessages = append(ErrorMessages, err.Error())
		return err
	}

	if isValid != true {
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

	response, err := db.ExecContext(ctx, UpdateUserCredQuery, Req.Password, Req.UserID)

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
