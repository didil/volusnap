package api

type providerSvcer interface {
	ListVolumes() ([]Volume, error)
}
