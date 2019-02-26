package api

import (
	"database/sql"
	"testing"

	"github.com/didil/volusnap/pkg/models"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/boil"
)

type SnapRuleTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (suite *SnapRuleTestSuite) SetupSuite() {
	db := bootstrapTests()
	suite.db = db
}
func TestSnapRuleTestSuite(t *testing.T) {
	suite.Run(t, new(SnapRuleTestSuite))
}

func (suite *SnapRuleTestSuite) Test_snapRuleService_List() {
	db := suite.db
	defer func() {
		models.SnapRules().DeleteAll(db)
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	acc := models.Account{UserID: user.ID, Provider: "DigitalOcean"}
	err = acc.Insert(db, boil.Infer())
	suite.NoError(err)

	snapRule1 := models.SnapRule{AccountID: acc.ID, VolumeID: "vol-1", Frequency: 2}
	err = snapRule1.Insert(db, boil.Infer())
	suite.NoError(err)

	snapRule2 := models.SnapRule{AccountID: acc.ID, VolumeID: "vol-2", Frequency: 4}
	err = snapRule2.Insert(db, boil.Infer())
	suite.NoError(err)

	snapRuleSvc := newSnapRuleService(db)
	snapRules, err := snapRuleSvc.List(acc.ID)
	suite.NoError(err)

	suite.Len(snapRules, 2)
}

func (suite *SnapRuleTestSuite) Test_snapRuleService_CreateOk() {
	db := suite.db
	defer func() {
		models.SnapRules().DeleteAll(db)
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	acc := models.Account{UserID: user.ID, Provider: "DigitalOcean"}
	err = acc.Insert(db, boil.Infer())
	suite.NoError(err)

	snapRuleSvc := newSnapRuleService(db)
	snapRuleID, err := snapRuleSvc.Create(acc.ID, 12, "vol-123", "my volume")
	suite.NoError(err)

	snapRule, err := models.FindSnapRule(db, snapRuleID)
	suite.NoError(err)

	suite.Equal(acc.ID, snapRule.AccountID)
	suite.Equal(12, snapRule.Frequency)
	suite.Equal("vol-123", snapRule.VolumeID)
	suite.Equal("my volume", snapRule.VolumeName)
}
