package api

import (
	"fmt"

	"github.com/didil/volusnap/pkg/models"
)

// ProviderSvcer is the interface Provider Services must implement
type ProviderSvcer interface {
	ListVolumes() ([]Volume, error)
	TakeSnapshot(snapRule *models.SnapRule) (string, error)
}

func getProviderService(account *models.Account) (ProviderSvcer, error) {
	factory := pRegistry.getProviderServiceFactory(account.Provider)
	if factory == nil {
		return nil, fmt.Errorf("could not get provider factory for %v", account.Provider)
	}

	providerSvc := factory.Build(account.Token)
	return providerSvc, nil
}
