package api

import (
	"fmt"

	"github.com/didil/volusnap/pkg/models"
)

type snapshooter interface {
	Take(account *models.Account, volumeID string) error
}

func newSnapshotTaker() snapshooter {
	return &snapshotTaker{}
}

type snapshotTaker struct {
}

func (shooter *snapshotTaker) Take(account *models.Account, volumeID string) error {
	providerSvc, err := getProviderService(account)
	if err != nil {
		return fmt.Errorf("could not get provider service: %v", err)
	}

	err = providerSvc.TakeSnapshot(volumeID)
	return err
}
