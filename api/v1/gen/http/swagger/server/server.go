// Code generated by goa v3.14.1, DO NOT EDIT.
//
// swagger HTTP server
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/v1/design

package server

import (
	"context"
	"net/http"

	swagger "github.com/tektoncd/hub/api/v1/gen/swagger"
	goahttp "goa.design/goa/v3/http"
	"goa.design/plugins/v3/cors"
)

// Server lists the swagger service endpoint HTTP handlers.
type Server struct {
	Mounts             []*MountPoint
	CORS               http.Handler
	DocsV1Openapi3JSON http.Handler
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the swagger service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *swagger.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
	fileSystemDocsV1Openapi3JSON http.FileSystem,
) *Server {
	if fileSystemDocsV1Openapi3JSON == nil {
		fileSystemDocsV1Openapi3JSON = http.Dir(".")
	}
	return &Server{
		Mounts: []*MountPoint{
			{"CORS", "OPTIONS", "/v1/schema/swagger.json"},
			{"docs/v1/openapi3.json", "GET", "/v1/schema/swagger.json"},
		},
		CORS:               NewCORSHandler(),
		DocsV1Openapi3JSON: http.FileServer(fileSystemDocsV1Openapi3JSON),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "swagger" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.CORS = m(s.CORS)
}

// MethodNames returns the methods served.
func (s *Server) MethodNames() []string { return swagger.MethodNames[:] }

// Mount configures the mux to serve the swagger endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountCORSHandler(mux, h.CORS)
	MountDocsV1Openapi3JSON(mux, goahttp.Replace("", "/docs/v1/openapi3.json", h.DocsV1Openapi3JSON))
}

// Mount configures the mux to serve the swagger endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountDocsV1Openapi3JSON configures the mux to serve GET request made to
// "/v1/schema/swagger.json".
func MountDocsV1Openapi3JSON(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/v1/schema/swagger.json", HandleSwaggerOrigin(h).ServeHTTP)
}

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service swagger.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleSwaggerOrigin(h)
	mux.Handle("OPTIONS", "/v1/schema/swagger.json", h.ServeHTTP)
}

// NewCORSHandler creates a HTTP handler which returns a simple 204 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
}

// HandleSwaggerOrigin applies the CORS response headers corresponding to the
// origin for the service swagger.
func HandleSwaggerOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "*") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET")
				w.WriteHeader(204)
				return
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
