package skiprouter_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/andrdru/openapi-mock/skiprouter"
)

func TestRegistrar_Add(t *testing.T) {
	type testEnv struct {
		router       *skiprouter.Registrar
		ctrl         *gomock.Controller
		routeHandler *Mockrouter
	}

	type testCase struct {
		name   string
		setup  func(te *testEnv) (pattern string)
		expect func() error
	}

	var newTestCase = func(t *testing.T) *testEnv {
		ctrl := gomock.NewController(t)

		routeHandler := NewMockrouter(ctrl)

		return &testEnv{
			routeHandler: routeHandler,
			router:       skiprouter.NewRegistrar(routeHandler),
		}
	}

	var handler http.HandlerFunc = func(http.ResponseWriter, *http.Request) {}

	var tests = []testCase{
		{
			name: "add success",
			setup: func(te *testEnv) (pattern string) {
				pattern = "/my/route"
				te.routeHandler.EXPECT().HandleFunc(pattern, gomock.Any())

				return pattern
			},
			expect: func() error {
				return nil
			},
		},
		{
			name: "err duplicate",
			setup: func(te *testEnv) (pattern string) {
				pattern = "/my/route"

				te.routeHandler.EXPECT().HandleFunc(pattern, gomock.Any())
				_ = te.router.Add(pattern, handler)

				return pattern
			},
			expect: func() error {
				return skiprouter.ErrDuplicatedRoute
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			te := newTestCase(tt)

			patternSet := test.setup(te)
			errExp := test.expect()

			errRet := te.router.Add(patternSet, handler)

			assert.ErrorIs(tt, errRet, errExp)
		})
	}
}
