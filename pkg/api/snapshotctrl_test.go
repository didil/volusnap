package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/didil/volusnap/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSnapshotSvc struct {
	mock.Mock
}

func (m *mockSnapshotSvc) List(snapRuleID int) (models.SnapshotSlice, error) {
	args := m.Called(snapRuleID)
	return args.Get(0).(models.SnapshotSlice), args.Error(1)
}

func (m *mockSnapshotSvc) ExistsFor(snapRuleID int, createdAfter time.Time) (bool, error) {
	args := m.Called(snapRuleID, createdAfter)
	return args.Bool(0), args.Error(1)
}

func (m *mockSnapshotSvc) Create(snapRuleID int) (int, error) {
	args := m.Called(snapRuleID)
	return args.Int(0), args.Error(1)
}
func Test_handleListSnapshotsOK(t *testing.T) {
	userID := 5
	token, err := signJWT(userID)
	assert.NoError(t, err)

	snapshotSvc := new(mockSnapshotSvc)
	accountSvc := new(mockAccountSvc)
	snapshotCtrl := newSnapshotController(snapshotSvc, accountSvc)

	accountID := 101
	account := &models.Account{ID: accountID, Provider: "test-provider", Token: "test-token"}

	accountSvc.On("GetForUser", userID, accountID).Return(account, nil)

	snapshots := models.SnapshotSlice{
		&models.Snapshot{ID: 28},
	}

	snapshotSvc.On("List", accountID).Return(snapshots, nil)

	r := buildRouter(&appController{snapshotCtrl: snapshotCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL+fmt.Sprintf("/api/v1/account/%v/snapshot/", accountID), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/JSON")

	var listResp listSnapshotsResp
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	assert.NoError(t, err)

	assert.ElementsMatch(t, listResp.Snapshots, snapshots)

	accountSvc.AssertExpectations(t)
	snapshotSvc.AssertExpectations(t)
}
