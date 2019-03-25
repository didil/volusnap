package api

import "sync"

var pRegistry = newProviderRegistry()

type providerServiceFactory interface {
	Build(token string) ProviderSvcer
}

func newProviderRegistry() *providerRegistry {
	reg := &providerRegistry{}
	reg.factories = make(map[string]providerServiceFactory)
	return reg
}

// singleton providers registry
type providerRegistry struct {
	mutex     sync.Mutex
	factories map[string]providerServiceFactory
}

// register a new provider
func (reg *providerRegistry) register(provider string, factory providerServiceFactory) {
	reg.mutex.Lock()
	defer reg.mutex.Unlock()
	reg.factories[provider] = factory
}

// check if provider valid
func (reg *providerRegistry) isValidProvider(provider string) bool {
	reg.mutex.Lock()
	defer reg.mutex.Unlock()
	_, ok := reg.factories[provider]
	return ok
}

// get provider service factory
func (reg *providerRegistry) getProviderServiceFactory(provider string) providerServiceFactory {
	reg.mutex.Lock()
	defer reg.mutex.Unlock()
	factory, ok := reg.factories[provider]
	if !ok {
		return nil
	}
	return factory
}
