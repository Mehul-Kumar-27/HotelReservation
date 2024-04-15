package api

import "github.com/gin-gonic/gin"

type Router interface {
	Routes() []Route
	Path() string
	AddToGroup(parentGroup *gin.RouterGroup) *gin.RouterGroup
}

type Route interface {
	Handlers() gin.HandlersChain
	Path() string
	Method() string
}

type AbstractRouter struct {
	Router

	routes []Route
	path   string
}

type AbstractRoute struct {
	Route

	path    string
	method  string
	handler gin.HandlersChain
}

func (r *AbstractRouter) Path() string {
	return r.path
}

func (r *AbstractRouter) Routes() []Route {
	return r.routes
}

func (r *AbstractRouter) AddToGroup(parentGroup *gin.RouterGroup) *gin.RouterGroup {

	childGroup := parentGroup.Group(r.Path())
	childGroup = r.addHandlers(childGroup)

	return parentGroup

}

func (r *AbstractRouter) addHandlers(parentGroup *gin.RouterGroup) *gin.RouterGroup {

	for _, route := range r.Routes() {
		parentGroup.Handle(route.Method(), route.Path(), route.Handlers()...)
	}
	return parentGroup
}

func NewRouter(path string, routes []Route) *AbstractRouter {
	return &AbstractRouter{
		path:   path,
		routes: routes,
	}
}

func (r *AbstractRoute) Handlers() gin.HandlersChain {
	return r.handler
}
func (r *AbstractRoute) Method() string {
	return r.method
}

func (r *AbstractRoute) Path() string {
	return r.path
}

func NewRoute(path, method string, handler ...gin.HandlerFunc) *AbstractRoute {
	return &AbstractRoute{
		path:    path,
		method:  method,
		handler: handler,
	}
}
