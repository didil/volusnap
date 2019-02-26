package api

import (
	"database/sql"
	"testing"
	"time"

	"github.com/volatiletech/null"

	"github.com/didil/volusnap/pkg/models"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/boil"
)

type SnapshotTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (suite *SnapshotTestSuite) SetupSuite() {
	db := bootstrapTests()
	suite.db = db
}
func TestSnapshotTestSuite(t *testing.T) {
	suite.Run(t, new(SnapshotTestSuite))
}

func (suite *SnapshotTestSuite) Test_snapshotService_List() {
	db := suite.db
	defer func() {
		models.Snapshots().DeleteAll(db)
		models.SnapRules().DeleteAll(db)
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	acc1 := models.Account{UserID: user.ID, Provider: "DigitalOcean"}
	err = acc1.Insert(db, boil.Infer())
	suite.NoError(err)

	acc2 := models.Account{UserID: user.ID, Provider: "DigitalOcean"}
	err = acc2.Insert(db, boil.Infer())
	suite.NoError(err)

	snapRule1 := models.SnapRule{AccountID: acc1.ID, VolumeID: "vol-1", Frequency: 2}
	err = snapRule1.Insert(db, boil.Infer())
	suite.NoError(err)

	snapRule2 := models.SnapRule{AccountID: acc2.ID, VolumeID: "vol-2", Frequency: 2}
	err = snapRule2.Insert(db, boil.Infer())
	suite.NoError(err)

	snapshot1 := models.Snapshot{SnapRuleID: snapRule1.ID}
	err = snapshot1.Insert(db, boil.Infer())
	suite.NoError(err)

	snapshot2 := models.Snapshot{SnapRuleID: snapRule1.ID}
	err = snapshot2.Insert(db, boil.Infer())
	suite.NoError(err)

	snapshot3 := models.Snapshot{SnapRuleID: snapRule2.ID}
	err = snapshot3.Insert(db, boil.Infer())
	suite.NoError(err)

	snapshotSvc := newSnapshotService(db)
	snapshots, err := snapshotSvc.List(acc1.ID)
	suite.NoError(err)

	suite.Len(snapshots, 2)
}

func (suite *SnapshotTestSuite) Test_snapshotService_ExistsFor() {
	db := suite.db
	defer func() {
		models.Snapshots().DeleteAll(db)
		models.SnapRules().DeleteAll(db)
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	acc1 := models.Account{UserID: user.ID, Provider: "DigitalOcean"}
	err = acc1.Insert(db, boil.Infer())
	suite.NoError(err)

	snapRule1 := models.SnapRule{AccountID: acc1.ID, VolumeID: "vol-1", Frequency: 2}
	err = snapRule1.Insert(db, boil.Infer())
	suite.NoError(err)
	createdAfter := time.Now().Add(-2 * time.Hour)

	snapshot1 := models.Snapshot{SnapRuleID: snapRule1.ID, CreatedAt: null.TimeFrom(time.Now().Add(-3 * time.Hour))}
	err = snapshot1.Insert(db, boil.Infer())
	suite.NoError(err)

	snapshotSvc := newSnapshotService(db)
	exists, err := snapshotSvc.ExistsFor(snapRule1.ID, createdAfter)
	suite.NoError(err)

	suite.Equal(false, exists)

	snapshot2 := models.Snapshot{SnapRuleID: snapRule1.ID, CreatedAt: null.TimeFrom(time.Now().Add(-1 * time.Hour))}
	err = snapshot2.Insert(db, boil.Infer())
	suite.NoError(err)

	snapshotSvc = newSnapshotService(db)
	exists, err = snapshotSvc.ExistsFor(snapRule1.ID, createdAfter)
	suite.NoError(err)

	suite.Equal(true, exists)
}

func (suite *SnapshotTestSuite) Test_snapshotService_Create() {
	db := suite.db
	defer func() {
		models.Snapshots().DeleteAll(db)
		models.SnapRules().DeleteAll(db)
		models.Accounts().DeleteAll(db)
		models.Users().DeleteAll(db)
	}()

	user := models.User{Email: "ex@example.com"}
	err := user.Insert(db, boil.Infer())
	suite.NoError(err)

	acc1 := models.Account{UserID: user.ID, Provider: "DigitalOcean"}
	err = acc1.Insert(db, boil.Infer())
	suite.NoError(err)

	snapRule1 := models.SnapRule{AccountID: acc1.ID, VolumeID: "vol-1", Frequency: 2}
	err = snapRule1.Insert(db, boil.Infer())
	suite.NoError(err)

	snapshotSvc := newSnapshotService(db)
	snapshotID, err := snapshotSvc.Create(snapRule1.ID)
	suite.NoError(err)

	snapshot, err := models.FindSnapshot(db, snapshotID)
	suite.NoError(err)

	suite.Equal(snapshot.ID, snapshotID)
	suite.Equal(snapshot.SnapRuleID, snapRule1.ID)
}
