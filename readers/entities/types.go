package entities

type (
	Response map[string]any

	Route struct {
		Pattern  string
		Response Response
	}

	Routes []Route
)
