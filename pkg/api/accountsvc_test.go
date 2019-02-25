package api

import (
	"database/sql"

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

func (suite *AuthTestSuite) Test_accountService_List() {
	db := suite.db
	defer func() {
		models.Accounts().DeleteAll(db)
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
