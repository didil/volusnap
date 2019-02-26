package api

import (
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/didil/volusnap/pkg/models"
)

type snapshotSvcer interface {
	List(snapRuleID int) (models.SnapshotSlice, error)
	ExistsFor(snapRuleID int, createdAfter time.Time) (bool, error)
	Create(snapRuleID int, providerSnapshotID string) (int, error)
}

func newSnapshotService(db *sql.DB) *snapshotService {
	return &snapshotService{db}
}

type snapshotService struct {
	db *sql.DB
}

func (svc *snapshotService) List(accountID int) (models.SnapshotSlice, error) {
	snapshots, err := models.Snapshots(
		qm.InnerJoin("snap_rules on snapshots.snap_rule_id = snap_rules.id"),
		qm.Where("snap_rules.account_id=?", accountID),
	).All(svc.db)
	return snapshots, err
}

func (svc *snapshotService) ExistsFor(snapRuleID int, createdAfter time.Time) (bool, error) {
	count, err := models.Snapshots(
		qm.Where("snap_rule_id=? AND created_at > ?", snapRuleID, createdAfter),
	).Count(svc.db)
	return count > 0, err
}

func (svc *snapshotService) Create(snapRuleID int, providerSnapshotID string) (int, error) {
	snapshot := models.Snapshot{
		SnapRuleID:         snapRuleID,
		ProviderSnapshotID: providerSnapshotID,
	}
	err := snapshot.Insert(svc.db, boil.Infer())

	if err != nil {
		return 0, err
	}

	return snapshot.ID, nil
}
