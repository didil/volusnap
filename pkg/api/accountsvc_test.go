package api

import (
	"database/sql"
	"testing"

	"github.com/didil/volusnap/pkg/models"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/boil"
)

type AccountTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (suite *AccountTestSuite) SetupSuite() {
	db := bootstrapTests()
	suite.db = db
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

func (suite *AuthTestSuite) Test_accountService_List() {
	db := suite.db
	defer func() {
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	acc1 := models.Account{UserID: user.ID, Provider: "DigitalOcean"}
	err = acc1.Insert(db, boil.Infer())
	suite.NoError(err)

	acc2 := models.Account{UserID: user.ID, Provider: "Scaleway"}
	err = acc2.Insert(db, boil.Infer())
	suite.NoError(err)

	accountSvc := newAccountService(db)
	accounts, err := accountSvc.List(user.ID)
	suite.NoError(err)

	suite.Len(accounts, 2)
}

func (suite *AuthTestSuite) Test_accountService_CreateOk() {
	db := suite.db
	defer func() {
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	accountSvc := newAccountService(db)
	accountID, err := accountSvc.Create(user.ID, "digital_ocean", "do 1", "token123")
	suite.NoError(err)

	account, err := models.FindAccount(db, accountID)
	suite.NoError(err)

	suite.Equal(user.ID, account.UserID)
	suite.Equal("digital_ocean", account.Provider)
	suite.Equal("do 1", account.Name)
	suite.Equal("token123", account.Token)
}

func (suite *AuthTestSuite) Test_accountService_CreateInvalidProvider() {
	db := suite.db
	defer func() {
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	accountSvc := newAccountService(db)
	_, err = accountSvc.Create(user.ID, "gammax_cloud", "do 1", "token123")
	suite.EqualError(err, "invalid provider: gammax_cloud")
}

func (suite *AuthTestSuite) Test_accountService_Get() {
	db := suite.db
	defer func() {
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	acc1 := models.Account{UserID: user.ID, Provider: "DigitalOcean"}
	err = acc1.Insert(db, boil.Infer())
	suite.NoError(err)

	acc2 := models.Account{UserID: user.ID, Provider: "Scaleway"}
	err = acc2.Insert(db, boil.Infer())
	suite.NoError(err)

	accountSvc := newAccountService(db)
	account, err := accountSvc.Get(user.ID, acc1.ID)
	suite.NoError(err)

	suite.Equal(account.Provider, acc1.Provider)
	suite.Equal(account.ID, acc1.ID)
}
