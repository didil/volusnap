package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/didil/volusnap/pkg/models"
)

func Test_snapshotTaker_Take(t *testing.T) {
	snapRule := &models.SnapRule{VolumeID: "my-vol-101"}
	providerSnapshotID := "my-snap-5"

	account := &models.Account{Provider: "test-provider", Token: "test-token"}

	providerSvc := new(mockProviderSvc)
	providerSvc.On("TakeSnapshot", snapRule).Return(providerSnapshotID, nil)

	pServiceFactory := new(mockProviderServiceFactory)
	pServiceFactory.On("Build", "test-token").Return(providerSvc)

	pRegistry.register("test-provider", pServiceFactory)

	snapshotTaker := newSnapshotTaker()
	myProviderSnapshotID, err := snapshotTaker.Take(account, snapRule)
	assert.NoError(t, err)

	assert.Equal(t, providerSnapshotID, myProviderSnapshotID)

	providerSvc.AssertExpectations(t)
	pServiceFactory.AssertExpectations(t)
}
