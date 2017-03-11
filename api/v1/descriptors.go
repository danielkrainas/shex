package v1

import (
	"net/http"
	//"regexp"

	"github.com/danielkrainas/gobag/api/describe"
	"github.com/danielkrainas/gobag/api/errcode"
)

var (
	versionHeaders = []describe.Parameter{
		{

			Name:        "Shex-Registry-Version",
			Type:        "string",
			Description: "The build version of the Shex registry server.",
			Format:      "<version>",
			Examples:    []string{"0.0.0-dev"},
		},
		{
			Name:        "Shex-Registry-Version",
			Type:        "string",
			Description: "The highest api version supported by the server.",
			Format:      "<version>",
			Examples:    []string{"1"},
		},
	}

	hostHeader = describe.Parameter{
		Name:        "Host",
		Type:        "string",
		Description: "",
		Format:      "<hostname>",
		Examples:    []string{"api.shexr.io"},
	}

	jsonContentLengthHeader = describe.Parameter{
		Name:        "Content-Length",
		Type:        "integer",
		Description: "Length of the JSON body.",
		Format:      "<length>",
	}
)

var API = struct {
	Routes []describe.Route `json:"routes"`
}{
	Routes: routeDescriptors,
}

var routeDescriptors = []describe.Route{
	{
		Name:        RouteNameBase,
		Path:        "/v1",
		Entity:      "Base",
		Description: "Base V1 API route, can be used for lightweight health and version check.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Check that the server supports the Shex V1 API.",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "The API implements the V1 protocol and is accessible.",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									jsonContentLengthHeader,
									versionHeaders...,
								},
							},
						},

						Failures: []describe.Response{
							{
								Description: "The API does not support the V1 protocol.",
								StatusCode:  http.StatusNotFound,
								Headers: []describe.Parameter{
									versionHeaders...,
								},
							},
						},
					},
				},
			},
		},
	},
}

var routeDescriptorsMap map[string]describe.Route

func init() {
	routeDescriptorsMap = make(map[string]describe.Route, len(routeDescriptors))
	for _, descriptor := range routeDescriptors {
		routeDescriptorsMap[descriptor.Name] = descriptor
	}
}
