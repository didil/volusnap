package api

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func bootstrapTests() *gorm.DB {
	err := loadConfig("../../config_test")
	if err != nil {
		panic(err)
	}

	db, err := openDB()
	if err != nil {
		panic(err)
	}

	err = autoMigrate(db)
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
		db.Delete(&User{})
	}()

	email := "example@example.com"
	password := "123456"

	err := db.Create(&User{Email: email}).Error
	suite.NoError(err)

	authSvc := newAuthService(db)
	_, err = authSvc.Signup(email, password)
	suite.EqualError(err, "a user already exists for the email: example@example.com")
}

func (suite *AuthTestSuite) Test_authService_SignupOK() {
	db := suite.db
	defer func() {
		db.Delete(&User{})
	}()

	email := "example@example.com"
	password := "123456"

	authSvc := newAuthService(db)
	id, err := authSvc.Signup(email, password)
	suite.NoError(err)

	user := User{}
	err = db.First(&user, id).Error
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
		db.Delete(&User{})
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
		db.Delete(&User{})
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
		db.Delete(&User{})
	}()

	email := "example@example.com"
	password := "123456"

	authSvc := newAuthService(db)

	id, err := authSvc.Signup(email, password)
	suite.NoError(err)

	token, err := authSvc.Login(email, password)
	suite.NoError(err)

	userID, err := parseJWT(token, []byte(viper.GetString("jwt.secret")))
	suite.NoError(err)

	suite.Equal(id, userID)
}
