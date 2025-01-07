package openapi_mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/andrdru/openapi-mock/readers/entities"
)

type (
	Mock struct {
		logger logger

		routes              entities.Routes
		routePatternToIndex map[string]int
		mockRoutesReader    routeReader

		mu sync.Mutex
	}

	routeReader interface {
		Read() (entities.Routes, error)
	}

	routeSkipper interface {
		Add(pattern string, handler http.HandlerFunc) error
	}

	logger interface {
		Warn(msg string, args ...any)
		Error(msg string, args ...any)
	}
)

const (
	ContentType = "application/json"
)

func NewMock(mockRoutesReader routeReader, logger logger) (*Mock, error) {
	var (
		err error
		ret = &Mock{
			logger:              logger,
			routePatternToIndex: make(map[string]int),
			mockRoutesReader:    mockRoutesReader,
		}
	)

	ret.routes, err = ret.mockRoutesReader.Read()
	if err != nil {
		return nil, fmt.Errorf("read routes: %w", err)
	}

	// todo reuse
	//ret.updateRoutes()

	for index, route := range ret.routes {
		ret.routePatternToIndex[route.Pattern] = index
	}

	return ret, nil
}

func (m *Mock) InitRoutes(routeSkipper routeSkipper) (err error) {
	for _, route := range m.routes {
		err = routeSkipper.Add(route.Pattern, m.handlerFunc(route.Pattern))
		if err != nil {
			m.logger.Warn("skip route register",
				"path", route.Pattern,
				"err", err.Error())
		}
	}

	return nil
}

func (m *Mock) handlerFunc(pattern string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		m.mu.Lock()
		m.writeAnswer(w, m.routes[m.routePatternToIndex[pattern]].Response)
		m.mu.Unlock()

		go m.updateRoutes()
	}
}

func (m *Mock) writeAnswer(w http.ResponseWriter, response entities.Response) {
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			m.logger.Error("failed to write response", "err", err.Error())
		}

		return
	}

	w.Header().Set("Content-Type", ContentType)

	_, err = w.Write(data)
	if err != nil {
		m.logger.Error("failed to write response", "err", err.Error())
	}
}

func (m *Mock) updateRoutes() {
	routes, err := m.mockRoutesReader.Read()
	if err != nil {
		m.logger.Error("failed to read routes", "err", err.Error())
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, route := range routes {
		index, ok := m.routePatternToIndex[route.Pattern]
		if !ok {
			continue
		}

		m.routes[index] = route
	}
}
