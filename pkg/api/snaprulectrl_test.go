package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didil/volusnap/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSnapRuleSvc struct {
	mock.Mock
}

func (m *mockSnapRuleSvc) List(accountID int) (models.SnapRuleSlice, error) {
	args := m.Called(accountID)
	return args.Get(0).(models.SnapRuleSlice), args.Error(1)
}

func (m *mockSnapRuleSvc) ListAll() (models.SnapRuleSlice, error) {
	args := m.Called()
	return args.Get(0).(models.SnapRuleSlice), args.Error(1)
}

func (m *mockSnapRuleSvc) Create(accountID int, frequency int, volumeID string, volumeName string, volumeRegion string) (int, error) {
	args := m.Called(accountID, frequency, volumeID, volumeName, volumeRegion)
	return args.Int(0), args.Error(1)
}
func Test_handleListSnapRulesOK(t *testing.T) {
	userID := 5
	token, err := signJWT(userID)
	assert.NoError(t, err)

	snapRuleSvc := new(mockSnapRuleSvc)
	accountSvc := new(mockAccountSvc)
	snapRuleCtrl := newSnapRuleController(snapRuleSvc, accountSvc)

	accountID := 101
	account := &models.Account{ID: accountID, Provider: "test-provider", Token: "test-token"}

	accountSvc.On("GetForUser", userID, accountID).Return(account, nil)

	snapRules := models.SnapRuleSlice{
		&models.SnapRule{ID: 15},
	}

	snapRuleSvc.On("List", accountID).Return(snapRules, nil)

	r := buildRouter(&appController{snapRuleCtrl: snapRuleCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL+fmt.Sprintf("/api/v1/account/%v/snaprule/", accountID), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/JSON")

	var listResp listSnapRulesResp
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	assert.NoError(t, err)

	assert.ElementsMatch(t, listResp.SnapRules, snapRules)

	accountSvc.AssertExpectations(t)
	snapRuleSvc.AssertExpectations(t)
}

func Test_handleCreateSnapRuleOK(t *testing.T) {
	userID := 5
	token, err := signJWT(userID)
	assert.NoError(t, err)

	snapRuleSvc := new(mockSnapRuleSvc)
	accountSvc := new(mockAccountSvc)
	snapRuleCtrl := newSnapRuleController(snapRuleSvc, accountSvc)

	accountID := 101
	account := &models.Account{ID: accountID, Provider: "test-provider", Token: "test-token"}

	accountSvc.On("GetForUser", userID, accountID).Return(account, nil)

	snapRuleID := 56

	frequency := 24
	volumeID := "vol-15"
	volumeName := "my-volu"
	volumeRegion := "lon1"

	snapRuleSvc.On("Create", accountID, frequency, volumeID, volumeName, volumeRegion).Return(snapRuleID, nil)

	r := buildRouter(&appController{snapRuleCtrl: snapRuleCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&createSnapRuleReq{Frequency: frequency, VolumeID: volumeID, VolumeName: volumeName, VolumeRegion: volumeRegion})

	req, err := http.NewRequest(http.MethodPost, s.URL+fmt.Sprintf("/api/v1/account/%v/snaprule/", accountID), &b)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/JSON")

	var createResp createSnapRuleResp
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	assert.NoError(t, err)

	assert.Equal(t, createResp.ID, snapRuleID)

	accountSvc.AssertExpectations(t)
}
