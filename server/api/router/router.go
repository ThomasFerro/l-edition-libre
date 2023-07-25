package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
)

type Route struct {
	Path        string
	Method      string
	Handler     middlewares.HandlerFuncReturningRequest
	Middlewares []middlewares.Middleware
}
type pathSegment struct {
	value string
}

func (segment pathSegment) isParameterized() bool {
	return strings.HasPrefix(segment.value, ":")
}

func (route Route) pathSegments() []pathSegment {
	segments := strings.Split(route.Path, "/")
	pathSegments := make([]pathSegment, 0)
	for _, segment := range segments {
		pathSegments = append(pathSegments, pathSegment{value: segment})
	}
	return pathSegments
}

func (route Route) firstParametrizedSegment() int {
	segments := route.pathSegments()
	for index, segment := range segments {
		if segment.isParameterized() {
			return index
		}
	}
	return -1
}

func (route Route) pathUntilFirstParametrizedSegment() string {
	path := ""
	for _, segment := range route.pathSegments() {
		if segment.isParameterized() {
			return path
		}
		path = fmt.Sprintf("%v%v/", path, segment.value)
	}
	return path
}

func (route Route) matchRequest(url, method string) bool {
	urlSegments := strings.Split(url, "/")
	pathSegments := route.pathSegments()
	if len(urlSegments) != len(pathSegments) {
		return false
	}
	for index, urlSegment := range urlSegments {
		if pathSegments[index].isParameterized() {
			continue
		}
		if pathSegments[index].value != urlSegment || method != route.Method {
			return false
		}
	}
	return true
}

func (route Route) addAllRouteParameters(ctx context.Context, url string) context.Context {
	urlSegments := strings.Split(url, "/")
	pathSegments := route.pathSegments()
	for index, urlSegment := range urlSegments {
		if !pathSegments[index].isParameterized() {
			continue
		}
		ctx = context.WithValue(ctx, fmt.Sprintf("URL_PARAM%v", pathSegments[index].value), urlSegment)
	}
	return ctx
}

func toHttpHandlerFunc(handler middlewares.HandlerFuncReturningRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

func (route Route) handlerWithMiddlewaresApplied() http.HandlerFunc {
	if len(route.Middlewares) == 0 {
		return toHttpHandlerFunc(route.Handler)
	}
	routeHandler := func(w http.ResponseWriter, r *http.Request) *http.Request {
		return route.Handler(w, r)
	}
	handler := route.Middlewares[0](routeHandler)
	for i := 1; i < len(route.Middlewares); i++ {
		handler = route.Middlewares[i](handler)
	}
	return toHttpHandlerFunc(handler)
}

func customHandlerFunc(routes []Route) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			requestUrl := r.URL.String()
			if !route.matchRequest(requestUrl, r.Method) {
				continue
			}
			r = r.WithContext(route.addAllRouteParameters(r.Context(), requestUrl))
			route.handlerWithMiddlewaresApplied()(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
	}
}

func HandleRoutes(serveMux *http.ServeMux, routes []Route) {
	routesByPathName := make(map[string][]Route, 0)
	for _, nextRoute := range routes {
		nextRoutePath := nextRoute.pathUntilFirstParametrizedSegment()
		firstParametrizedSegment := nextRoute.firstParametrizedSegment()
		if firstParametrizedSegment == -1 {
			nextRoutePath = nextRoutePath[:len(nextRoutePath)-1]
		}

		if _, alreadyExist := routesByPathName[nextRoutePath]; alreadyExist {
			routesByPathName[nextRoutePath] = append(routesByPathName[nextRoutePath], nextRoute)
		} else {
			routesByPathName[nextRoutePath] = []Route{nextRoute}
		}
	}
	for path, routes := range routesByPathName {
		serveMux.HandleFunc(path, customHandlerFunc(routes))
	}
}
