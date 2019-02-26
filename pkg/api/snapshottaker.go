package api

import (
	"fmt"

	"github.com/didil/volusnap/pkg/models"
)

type snapshooter interface {
	Take(account *models.Account, snapRule *models.SnapRule) (string, error)
}

func newSnapshotTaker() snapshooter {
	return &snapshotTaker{}
}

type snapshotTaker struct {
}

func (shooter *snapshotTaker) Take(account *models.Account, snapRule *models.SnapRule) (string, error) {
	providerSvc, err := getProviderService(account)
	if err != nil {
		return "", fmt.Errorf("could not get provider service: %v", err)
	}

	providerSnapshotID, err := providerSvc.TakeSnapshot(snapRule)
	return providerSnapshotID, err
}
