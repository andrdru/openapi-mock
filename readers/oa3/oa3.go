package oa3

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/url"
	"path"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/andrdru/openapi-mock/readers/entities"
)

type Reader struct {
	oa3t   *openapi3.T
	prefix string

	maxDepth              int64
	contentType           string
	arrayItemsDisplay     int64
	randomFillNonRequired bool
}

const (
	contentType       = "application/json"
	masDepth          = 3
	arrayItemsDisplay = 3
)

var (
	ErrNoSchema = errors.New("no schema detected")

	randInt = rand.Int

	randomFillNonRequired = new(struct{})
)

func NewReader(filePath string, options ...Option) (*Reader, error) {
	var (
		err error
		ret = &Reader{}
	)

	ret.setOptions(options...)

	ret.oa3t, err = openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("loading swagger json failed: %w", err)
	}

	if len(ret.oa3t.Servers) > 0 {
		var serverUrl *url.URL

		serverUrl, err = url.Parse(ret.oa3t.Servers[0].URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse server url: %w", err)
		}

		ret.prefix = serverUrl.Path
	}

	return ret, nil
}

func (r *Reader) setOptions(options ...Option) {
	var args = &Options{
		maxDepth:              masDepth,
		contentType:           contentType,
		arrayItemsDisplay:     arrayItemsDisplay,
		randomFillNonRequired: randomFillNonRequired,
	}

	var opt Option
	for _, opt = range options {
		opt(args)
	}

	r.maxDepth = args.maxDepth
	r.contentType = args.contentType
	r.arrayItemsDisplay = args.arrayItemsDisplay

	if args.randomFillNonRequired != nil {
		r.randomFillNonRequired = true
	}
}

func (r *Reader) Read() (routes entities.Routes, err error) {
	operations := r.getOperations()

	for operationPath, operation := range operations {
		route := entities.Route{
			Pattern:  operationPath,
			Response: r.getResponse(operation),
		}

		routes = append(routes, route)

	}

	return routes, nil
}

func (r *Reader) getOperations() map[string]*openapi3.Operation {
	ret := make(map[string]*openapi3.Operation)

	pathMap := r.oa3t.Paths.Map()
	for itemPath, item := range pathMap {
		operations := item.Operations()

		for httpMethod := range operations {
			ret[fmt.Sprintf("%s %s", httpMethod, path.Join(r.prefix, itemPath))] = operations[httpMethod]
		}
	}

	return ret
}

func (r *Reader) getResponse(operation *openapi3.Operation) entities.Response {
	responsesMap := operation.Responses.Map()

	schema, err := r.getResponseSchema(responsesMap)
	if err != nil {
		return entities.Response{}
	}

	return r.getProperties(schema.Properties, schema.Required, 0)
}

// getResponseSchema get response openapi schema
func (r *Reader) getResponseSchema(m map[string]*openapi3.ResponseRef) (schema *openapi3.Schema, err error) {
	key := getResponseKey(m)

	if _, ok := m[key]; !ok {
		return nil, fmt.Errorf("no key in responses %s: %w", key, ErrNoSchema)
	}

	if m[key].Value == nil {
		return nil, fmt.Errorf("no value in responses %s: %w", key, ErrNoSchema)
	}

	if _, ok := m[key].Value.Content[r.contentType]; !ok {
		return nil, fmt.Errorf("no content %s in response %s: %w", ContentType, key, ErrNoSchema)
	}

	if m[key].Value.Content[r.contentType].Schema == nil {
		return nil, fmt.Errorf("no schema in response %s: %w", key, ErrNoSchema)
	}

	if m[key].Value.Content[r.contentType].Schema.Value == nil {
		return nil, fmt.Errorf("no value in schema %s: %w", key, ErrNoSchema)
	}

	return m[key].Value.Content[r.contentType].Schema.Value, nil
}

func (r *Reader) getProperties(properties openapi3.Schemas, required []string, currentDepth int64) entities.Response {
	resp := make(map[string]any)

	for name, property := range properties {
		if property.Value.Type == nil || len(*property.Value.Type) == 0 {
			continue
		}

		valueType := (*property.Value.Type)[0]

		switch valueType {
		case "object":
			if currentDepth+1 >= r.maxDepth {
				resp[name] = nil

				continue
			}

			resp[name] = r.getProperties(
				property.Value.Properties,
				property.Value.Required,
				currentDepth+1,
			)

		case "array":
			arr := make([]any, 0, r.arrayItemsDisplay)
			for range r.arrayItemsDisplay {
				arr = append(arr,
					r.getProperties(
						property.Value.Items.Value.Properties,
						property.Value.Items.Value.Required,
						currentDepth,
					))
			}

			resp[name] = arr

		default:
			if !r.fillWithValue(required, name) {
				resp[name] = nil

				continue
			}

			resp[name] = getValue(property)
		}
	}

	return resp
}

func getValue(property *openapi3.SchemaRef) any {
	if property.Value == nil {
		//todo log error maybe
		return nil
	}

	return property.Value.Example
}

// getResponseKey get first key as response key
func getResponseKey(m map[string]*openapi3.ResponseRef) (key string) {
	for k := range m {
		return k
	}

	return ""
}

func (r *Reader) fillWithValue(required []string, name string) bool {
	if slices.Contains(required, name) {
		return true
	}

	if !r.randomFillNonRequired {
		return false
	}

	return randInt()%2 == 0
}
