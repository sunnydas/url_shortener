package routes

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"strings"
)

type RouteRegistrar struct {
	Router      *mux.Router
	routeConfig []*mux.Route
}

func (r *RouteRegistrar) RegisterRoute() *RouteRegistrar {
	r.Router.StrictSlash(false)
	r.routeConfig = []*mux.Route{r.Router.NewRoute(), r.Router.NewRoute()}
	return r
}

func (r *RouteRegistrar) Path(path string) *RouteRegistrar {
	trimmed := strings.TrimSuffix(path, "/")
	r.routeConfig[0].Path(trimmed)
	r.routeConfig[1].Path(trimmed + "/")

	return r
}

func (r *RouteRegistrar) Method(method string) *RouteRegistrar {
	for _, route := range r.routeConfig {
		route.Methods(method)
	}

	return r
}

func (r *RouteRegistrar) Query(params ...string) *RouteRegistrar {
	for _, route := range r.routeConfig {
		route.Queries(params...)
	}

	return r
}

func (r *RouteRegistrar) Handler(handler *negroni.Negroni) *RouteRegistrar {
	for _, route := range r.routeConfig {
		route.Handler(handler)
	}

	return r
}
