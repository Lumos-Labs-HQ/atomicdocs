package registry

import (
	"strconv"
	"sync"
	"github.com/yourusername/atomicdocs/internal/types"
)

type Registry struct {
	mu    sync.RWMutex
	apps  map[string][]types.RouteInfo
}

func New() *Registry {
	return &Registry{apps: make(map[string][]types.RouteInfo)}
}

func (r *Registry) RegisterApp(port int, routes []types.RouteInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.apps[strconv.Itoa(port)] = routes
}

func (r *Registry) GetByPort(port string) []types.RouteInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if routes, ok := r.apps[port]; ok {
		return routes
	}
	return []types.RouteInfo{}
}
