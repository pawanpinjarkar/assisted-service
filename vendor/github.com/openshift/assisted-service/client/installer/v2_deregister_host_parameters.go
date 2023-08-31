// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewV2DeregisterHostParams creates a new V2DeregisterHostParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewV2DeregisterHostParams() *V2DeregisterHostParams {
	return &V2DeregisterHostParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewV2DeregisterHostParamsWithTimeout creates a new V2DeregisterHostParams object
// with the ability to set a timeout on a request.
func NewV2DeregisterHostParamsWithTimeout(timeout time.Duration) *V2DeregisterHostParams {
	return &V2DeregisterHostParams{
		timeout: timeout,
	}
}

// NewV2DeregisterHostParamsWithContext creates a new V2DeregisterHostParams object
// with the ability to set a context for a request.
func NewV2DeregisterHostParamsWithContext(ctx context.Context) *V2DeregisterHostParams {
	return &V2DeregisterHostParams{
		Context: ctx,
	}
}

// NewV2DeregisterHostParamsWithHTTPClient creates a new V2DeregisterHostParams object
// with the ability to set a custom HTTPClient for a request.
func NewV2DeregisterHostParamsWithHTTPClient(client *http.Client) *V2DeregisterHostParams {
	return &V2DeregisterHostParams{
		HTTPClient: client,
	}
}

/*
V2DeregisterHostParams contains all the parameters to send to the API endpoint

	for the v2 deregister host operation.

	Typically these are written to a http.Request.
*/
type V2DeregisterHostParams struct {

	/* HostID.

	   The host that should be deregistered.

	   Format: uuid
	*/
	HostID strfmt.UUID

	/* InfraEnvID.

	   The infra-env of the host that should be deregistered.

	   Format: uuid
	*/
	InfraEnvID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the v2 deregister host params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V2DeregisterHostParams) WithDefaults() *V2DeregisterHostParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the v2 deregister host params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V2DeregisterHostParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the v2 deregister host params
func (o *V2DeregisterHostParams) WithTimeout(timeout time.Duration) *V2DeregisterHostParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v2 deregister host params
func (o *V2DeregisterHostParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v2 deregister host params
func (o *V2DeregisterHostParams) WithContext(ctx context.Context) *V2DeregisterHostParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v2 deregister host params
func (o *V2DeregisterHostParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v2 deregister host params
func (o *V2DeregisterHostParams) WithHTTPClient(client *http.Client) *V2DeregisterHostParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v2 deregister host params
func (o *V2DeregisterHostParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithHostID adds the hostID to the v2 deregister host params
func (o *V2DeregisterHostParams) WithHostID(hostID strfmt.UUID) *V2DeregisterHostParams {
	o.SetHostID(hostID)
	return o
}

// SetHostID adds the hostId to the v2 deregister host params
func (o *V2DeregisterHostParams) SetHostID(hostID strfmt.UUID) {
	o.HostID = hostID
}

// WithInfraEnvID adds the infraEnvID to the v2 deregister host params
func (o *V2DeregisterHostParams) WithInfraEnvID(infraEnvID strfmt.UUID) *V2DeregisterHostParams {
	o.SetInfraEnvID(infraEnvID)
	return o
}

// SetInfraEnvID adds the infraEnvId to the v2 deregister host params
func (o *V2DeregisterHostParams) SetInfraEnvID(infraEnvID strfmt.UUID) {
	o.InfraEnvID = infraEnvID
}

// WriteToRequest writes these params to a swagger request
func (o *V2DeregisterHostParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param host_id
	if err := r.SetPathParam("host_id", o.HostID.String()); err != nil {
		return err
	}

	// path param infra_env_id
	if err := r.SetPathParam("infra_env_id", o.InfraEnvID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}