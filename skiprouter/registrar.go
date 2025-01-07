package skiprouter

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

//go:generate mockgen -source=registrar.go -destination=registrar_mocks_test.go -package=skiprouter_test
type (
	router interface {
		HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	}

	Registrar struct {
		router        router
		existedRoutes map[string]struct{}

		mu sync.Mutex
	}
)

var (
	ErrDuplicatedRoute = errors.New("duplicated route")
)

func NewRegistrar(routeHandler router) *Registrar {
	return &Registrar{
		router:        routeHandler,
		existedRoutes: make(map[string]struct{}),
	}
}

func (r *Registrar) Add(pattern string, handler http.HandlerFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.existedRoutes[pattern]; ok {
		return fmt.Errorf("%s: %w", pattern, ErrDuplicatedRoute)
	}

	r.existedRoutes[pattern] = struct{}{}
	r.router.HandleFunc(pattern, handler)

	return nil
}
