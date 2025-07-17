package provider

import "sync"

type Registry struct {
	mutex     sync.RWMutex
	providers []Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make([]Provider, 0),
	}
}

func (r *Registry) Register(provider Provider) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.providers = append(r.providers, provider)
}

func (r *Registry) Unregister(provider Provider) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	pvds := r.providers
	for i, p := range pvds {
		if p == provider {
			r.providers = append(r.providers[:i], r.providers[i+1:]...)
			break
		}
	}
}

func (r *Registry) FindByQuery(query string) (Provider, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pvrs := r.providers
	for _, pvr := range pvrs {
		if pvr.CanHandle(query) {
			return pvr, true
		}
	}

	return nil, false
}

func (r *Registry) Providers() []Provider {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pvrs := make([]Provider, len(r.providers))
	copy(pvrs, r.providers)

	return pvrs
}
