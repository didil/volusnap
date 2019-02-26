package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/didil/volusnap/pkg/models"
)

func Test_snapshotTaker_Take(t *testing.T) {
	volumeID := "my-vol-101"

	account := &models.Account{Provider: "test-provider", Token: "test-token"}

	providerSvc := new(mockProviderSvc)
	providerSvc.On("TakeSnapshot", volumeID).Return(nil)

	pServiceFactory := new(mockProviderServiceFactory)
	pServiceFactory.On("Build", "test-token").Return(providerSvc)

	pRegistry.register("test-provider", pServiceFactory)

	snapshotTaker := newSnapshotTaker()
	err := snapshotTaker.Take(account, volumeID)
	assert.NoError(t, err)

	providerSvc.AssertExpectations(t)
	pServiceFactory.AssertExpectations(t)
}
