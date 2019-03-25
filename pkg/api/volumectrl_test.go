package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didil/volusnap/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProviderSvc struct {
	mock.Mock
}

func (m *mockProviderSvc) ListVolumes() ([]Volume, error) {
	args := m.Called()
	return args.Get(0).([]Volume), args.Error(1)
}

func (m *mockProviderSvc) TakeSnapshot(snapRule *models.SnapRule) (string, error) {
	args := m.Called(snapRule)
	return args.String(0), args.Error(1)
}

type mockProviderServiceFactory struct {
	mock.Mock
}

func (m *mockProviderServiceFactory) Build(token string) ProviderSvcer {
	args := m.Called(token)
	return args.Get(0).(ProviderSvcer)
}

func Test_handleListVolumesOK(t *testing.T) {
	userID := 5
	token, err := signJWT(userID)
	assert.NoError(t, err)

	accountSvc := new(mockAccountSvc)
	volumeCtrl := newVolumeController(accountSvc)

	accountID := 101
	account := &models.Account{Provider: "test-provider", Token: "test-token"}

	accountSvc.On("GetForUser", userID, accountID).Return(account, nil)

	volumes := []Volume{
		Volume{ID: "x-1", Name: "volume-name", Size: 5},
	}

	providerSvc := new(mockProviderSvc)
	providerSvc.On("ListVolumes").Return(volumes, nil)

	pServiceFactory := new(mockProviderServiceFactory)
	pServiceFactory.On("Build", "test-token").Return(providerSvc)

	pRegistry.register("test-provider", pServiceFactory)

	r := buildRouter(&appController{volumeCtrl: volumeCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL+fmt.Sprintf("/api/v1/account/%v/volume/", accountID), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/JSON")

	var listResp listVolumesResp
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	assert.NoError(t, err)

	assert.ElementsMatch(t, listResp.Volumes, volumes)

	accountSvc.AssertExpectations(t)
	pServiceFactory.AssertExpectations(t)
	providerSvc.AssertExpectations(t)
}
