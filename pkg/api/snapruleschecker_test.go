package api

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/didil/volusnap/pkg/models"

	"github.com/stretchr/testify/assert"
)

type mockSnapshotTaker struct {
	mock.Mock
}

func (m *mockSnapshotTaker) Take(account *models.Account, volumeID string) error {
	args := m.Called(account, volumeID)
	return args.Error(0)
}

func Test_snapRulesChecker_checkAll(t *testing.T) {
	snapRuleSvc := new(mockSnapRuleSvc)
	snapshotSvc := new(mockSnapshotSvc)
	accountSvc := new(mockAccountSvc)
	shooter := new(mockSnapshotTaker)

	account := &models.Account{ID: 1}

	snapshotID := 66

	snapRules := models.SnapRuleSlice{
		&models.SnapRule{ID: 5, Frequency: 12, AccountID: account.ID, VolumeID: "volum-1"},
		&models.SnapRule{ID: 8, Frequency: 24, AccountID: account.ID, VolumeID: "volum-2"},
	}
	snapRuleSvc.On("ListAll").Return(snapRules, nil)

	snapshotSvc.On("ExistsFor", snapRules[0].ID, mock.AnythingOfType("time.Time")).Return(true, nil)
	snapshotSvc.On("ExistsFor", snapRules[1].ID, mock.AnythingOfType("time.Time")).Return(false, nil)

	snapshotSvc.On("Create", snapRules[1].ID).Return(snapshotID, nil)

	accountSvc.On("Get", account.ID).Return(account, nil)

	shooter.On("Take", account, snapRules[1].VolumeID).Return(nil)

	checker := newSnapRulesChecker(snapRuleSvc, snapshotSvc, accountSvc, shooter)

	err := checker.checkAll()
	assert.NoError(t, err)

	shooter.AssertExpectations(t)
	snapshotSvc.AssertExpectations(t)
	snapRuleSvc.AssertExpectations(t)
	accountSvc.AssertExpectations(t)
}
