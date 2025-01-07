# Mock openapi3 library written in golang

Serve implemented routes in the usual way  
Serve mocks from openapi file

## Usage

### Minimal example
```go
// any http router
router := http.NewServeMux()

// custom router registrar, avoid duplicated routes
routeRegistrar := skiprouter.NewRegistrar(router)

// register implemented routes
_ = routeRegistrar.Add("GET /hello", func (writer http.ResponseWriter, request *http.Request) {
  _, _ = writer.Write([]byte("world"))
})

// openapi file reader and mock handler
swaggerReader, _ := oa3.NewReader("swagger.yaml")
apiMock, _ := mock.NewMock(swaggerReader, slog.Default())

// register unimplemented routes to mock response
_ = apiMock.InitRoutes(routeRegistrar)

// any http server
httpServer := &http.Server{Addr: ":8080", Handler: router}
_ = httpServer.ListenAndServe()

```

### Run it
see [./examples](./examples)
