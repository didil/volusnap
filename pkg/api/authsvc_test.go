package api

import (
	"database/sql"
	"testing"

	"github.com/didil/volusnap/pkg/models"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/boil"
)

type AuthTestSuite struct {
	suite.Suite
	db *sql.DB
}

func bootstrapTests() *sql.DB {
	err := loadConfig("../../config_test")
	if err != nil {
		panic(err)
	}

	db, err := openDB()
	if err != nil {
		panic(err)
	}

	return db
}

func (suite *AuthTestSuite) SetupSuite() {
	db := bootstrapTests()
	suite.db = db
}

func (suite *AuthTestSuite) Test_authService_SignupExisting() {
	db := suite.db
	defer func() {
		models.Users().DeleteAll(db)
	}()

	email := "example@example.com"
	password := "123456"

	user := models.User{Email: email}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	authSvc := newAuthService(db)
	_, err = authSvc.Signup(email, password)
	suite.EqualError(err, "a user already exists for the email: example@example.com")
}

func (suite *AuthTestSuite) Test_authService_SignupOK() {
	db := suite.db
	defer func() {
		models.Users().DeleteAll(db)
	}()

	email := "example@example.com"
	password := "123456"

	authSvc := newAuthService(db)
	id, err := authSvc.Signup(email, password)
	suite.NoError(err)

	user, err := models.FindUser(db, id)
	suite.NoError(err)

	suite.Equal(email, user.Email)
	suite.Len(user.Password, 60)
	suite.NoError(comparePasswords([]byte(user.Password), []byte(password)))
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (suite *AuthTestSuite) Test_authService_LoginNonExisting() {
	db := suite.db
	defer func() {
		models.Users().DeleteAll(db)
	}()

	email := "example@example.com"
	password := "123456"

	authSvc := newAuthService(db)
	_, err := authSvc.Login(email, password)
	suite.EqualError(err, "user not found for the email: example@example.com")
}

func (suite *AuthTestSuite) Test_authService_LoginInvalid() {
	db := suite.db
	defer func() {
		models.Users().DeleteAll(db)
	}()

	email := "example@example.com"
	password := "123456"

	authSvc := newAuthService(db)

	_, err := authSvc.Signup(email, password)
	suite.NoError(err)

	_, err = authSvc.Login(email, password+"x")
	suite.EqualError(err, "password invalid")
}

func (suite *AuthTestSuite) Test_authService_LoginOK() {
	db := suite.db
	defer func() {
		models.Users().DeleteAll(db)
	}()

	email := "example@example.com"
	password := "123456"

	authSvc := newAuthService(db)

	id, err := authSvc.Signup(email, password)
	suite.NoError(err)

	token, err := authSvc.Login(email, password)
	suite.NoError(err)

	userID, err := parseJWT(token)
	suite.NoError(err)

	suite.Equal(id, userID)
}
