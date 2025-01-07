package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	mock "github.com/andrdru/openapi-mock"
	"github.com/andrdru/openapi-mock/readers/oa3"
	"github.com/andrdru/openapi-mock/skiprouter"
)

func main() {
	log := slog.Default()

	log.Info("server starting")

	swaggerReader, err := oa3.NewReader("swagger.yaml",
		oa3.MaxDepth(3),
		oa3.ContentType("application/json"),
		oa3.RandomFillNonRequired(true),
		oa3.ArrayItemsDisplay(3),
	)
	if err != nil {
		log.Error("loading swagger json failed", "error", err.Error())
		os.Exit(1)
	}

	router := http.NewServeMux()

	routeRegistrar := skiprouter.NewRegistrar(router)

	_ = routeRegistrar.Add("/hello", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("world"))
	})

	// comment to check mock override
	//*
	_ = routeRegistrar.Add("GET /api/v1/users/profile2", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"replaced":"yes"}`))
	})
	// */

	apiMock, err := mock.NewMock(swaggerReader, slog.Default())
	if err != nil {
		log.Error("init mock", "error", err.Error())
		os.Exit(1)
	}

	err = apiMock.InitRoutes(routeRegistrar)
	if err != nil {
		log.Error("init routes", "error", err.Error())
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	log.Info("server started")

	if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("serve http", "error", err.Error())
	}
}
