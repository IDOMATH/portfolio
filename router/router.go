package router

import (
	"context"
	"net/http"
	"regexp"
)

//////////////////////////
// Router
//////////////////////////

type Router struct {
	routes []RouteEntry
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, entry := range rtr.routes {
		params := entry.Match(req)
		if params == nil {
			continue
		}

		// Create new request with params stored in context
		ctx := context.WithValue(req.Context(), "params", params)
		entry.Handler.ServeHTTP(w, req.WithContext(ctx))
		return
	}

	// No match found
	http.NotFound(w, req)
}

func (rtr *Router) Route(method, path string, handlerFunc http.HandlerFunc) {
	exactPath := regexp.MustCompile("^" + path + "$")
	e := RouteEntry{
		Method:  method,
		Path:    exactPath,
		Handler: handlerFunc,
	}
	rtr.routes = append(rtr.routes, e)
}

//////////////////////////
// RouteEntry
//////////////////////////

type RouteEntry struct {
	Path    *regexp.Regexp
	Method  string
	Handler http.HandlerFunc
}

func (route *RouteEntry) Match(req *http.Request) map[string]string {
	// Methods do not match
	if req.Method != route.Method {
		return nil
	}

	// Paths do not match
	match := route.Path.FindStringSubmatch(req.URL.Path)
	if match == nil {
		return nil
	}

	// Create map to store URL parameters
	params := make(map[string]string)
	groupNames := route.Path.SubexpNames()
	for i, group := range match {
		params[groupNames[i]] = group
	}

	return params
}

////////////////////
// Helpers
////////////////////

func UrlParam(req *http.Request, name string) string {
	ctx := req.Context()

	// ctx.Value returns an interface{}, so we need to cast it into a map
	// which is what we use to store parameters
	params := ctx.Value("params").(map[string]string)
	return params[name]
}
