package api

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type AccountTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *AccountTestSuite) SetupSuite() {
	db := bootstrapTests()
	suite.db = db
}

func (suite *AuthTestSuite) Test_accountService_List() {
	db := suite.db
	defer func() {
		db.Delete(&Account{})
	}()

	userID := uint(5)

	err := db.Create(&Account{UserID: userID, Provider: "DigitalOcean"}).Error
	suite.NoError(err)

	err = db.Create(&Account{UserID: userID, Provider: "Scaleway"}).Error
	suite.NoError(err)

	accountSvc := newAccountService(db)
	accounts, err := accountSvc.List(userID)
	suite.NoError(err)

	suite.Len(accounts, 2)
}
