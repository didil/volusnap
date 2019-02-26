package api

import (
	"database/sql"
	"fmt"

	"github.com/didil/volusnap/pkg/models"
	"github.com/volatiletech/sqlboiler/boil"
)

type snapRuleSvcer interface {
	List(accountID int) (models.SnapRuleSlice, error)
	ListAll() (models.SnapRuleSlice, error)
	Create(accountID int, frequency int, volumeID string, volumeName string, volumeRegion string) (int, error)
}

func newSnapRuleService(db *sql.DB) *snapRuleService {
	return &snapRuleService{db}
}

type snapRuleService struct {
	db *sql.DB
}

func (svc *snapRuleService) List(accountID int) (models.SnapRuleSlice, error) {
	snapRules, err := models.SnapRules(models.SnapRuleWhere.AccountID.EQ(accountID)).All(svc.db)
	return snapRules, err
}

func (svc *snapRuleService) ListAll() (models.SnapRuleSlice, error) {
	snapRules, err := models.SnapRules().All(svc.db)
	return snapRules, err
}

func (svc *snapRuleService) Create(accountID int, frequency int, volumeID string, volumeName string, volumeRegion string) (int, error) {
	if frequency == 0 {
		return 0, fmt.Errorf("frequency must be > 0")
	}
	if volumeID == "" {
		return 0, fmt.Errorf("empty volumeID")
	}

	snapRule := models.SnapRule{
		AccountID:    accountID,
		Frequency:    frequency,
		VolumeID:     volumeID,
		VolumeName:   volumeName,
		VolumeRegion: volumeRegion,
	}
	err := snapRule.Insert(svc.db, boil.Infer())

	if err != nil {
		return 0, err
	}

	return snapRule.ID, nil
}
