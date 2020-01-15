package actions

import (
	"reflect"

	"github.com/gobuffalo/buffalo"
)

// RequestMethod defines a HTTP Method so it describes the Functions from
// the 'App' type in the 'buffalo' Package.
// type RequestMethod

// RouteResource describes one route resource, it holds Informations – for
// internally use – about the route and the function.
type RouteResource struct {
	// Route contains the Route string
	Route string
	// Method describes which method is used for the Request
	// NOTE: Only add Functions which are available in the buffalo 'App' type
	// Method func(string, http.Handler) *buffalo.RouteInfo
	Method string
	// Handler contains the function which will handle the request
	Handler func(c buffalo.Context) error
	//
}

func (r *RouteResource) AddRoute(app *buffalo.App) *buffalo.RouteInfo {
	_method := reflect.ValueOf(app).MethodByName(r.Method)
	// get Method
	_cMethod := _method.Interface().(func(string, buffalo.Handler) *buffalo.RouteInfo)

	// add the new Route
	return _cMethod(r.Route, r.Handler)
}
