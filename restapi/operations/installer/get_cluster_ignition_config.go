// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetClusterIgnitionConfigHandlerFunc turns a function with the right signature into a get cluster ignition config handler
type GetClusterIgnitionConfigHandlerFunc func(GetClusterIgnitionConfigParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetClusterIgnitionConfigHandlerFunc) Handle(params GetClusterIgnitionConfigParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetClusterIgnitionConfigHandler interface for that can handle valid get cluster ignition config params
type GetClusterIgnitionConfigHandler interface {
	Handle(GetClusterIgnitionConfigParams, interface{}) middleware.Responder
}

// NewGetClusterIgnitionConfig creates a new http.Handler for the get cluster ignition config operation
func NewGetClusterIgnitionConfig(ctx *middleware.Context, handler GetClusterIgnitionConfigHandler) *GetClusterIgnitionConfig {
	return &GetClusterIgnitionConfig{Context: ctx, Handler: handler}
}

/*GetClusterIgnitionConfig swagger:route GET /clusters/{cluster_id}/ignition-config installer getClusterIgnitionConfig

Get the cluster ignition config

*/
type GetClusterIgnitionConfig struct {
	Context *middleware.Context
	Handler GetClusterIgnitionConfigHandler
}

func (o *GetClusterIgnitionConfig) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetClusterIgnitionConfigParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
