// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetTasksCountParams creates a new GetTasksCountParams object
// with the default values initialized.
func NewGetTasksCountParams() *GetTasksCountParams {
	var ()
	return &GetTasksCountParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetTasksCountParamsWithTimeout creates a new GetTasksCountParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetTasksCountParamsWithTimeout(timeout time.Duration) *GetTasksCountParams {
	var ()
	return &GetTasksCountParams{

		timeout: timeout,
	}
}

// NewGetTasksCountParamsWithContext creates a new GetTasksCountParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetTasksCountParamsWithContext(ctx context.Context) *GetTasksCountParams {
	var ()
	return &GetTasksCountParams{

		Context: ctx,
	}
}

// NewGetTasksCountParamsWithHTTPClient creates a new GetTasksCountParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetTasksCountParamsWithHTTPClient(client *http.Client) *GetTasksCountParams {
	var ()
	return &GetTasksCountParams{
		HTTPClient: client,
	}
}

/*GetTasksCountParams contains all the parameters to send to the API endpoint
for the get tasks count operation typically these are written to a http.Request
*/
type GetTasksCountParams struct {

	/*Repository
	  optional; repository for filtering

	*/
	Repository *string
	/*Sha
	  optional; sha for filtering

	*/
	Sha *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get tasks count params
func (o *GetTasksCountParams) WithTimeout(timeout time.Duration) *GetTasksCountParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get tasks count params
func (o *GetTasksCountParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get tasks count params
func (o *GetTasksCountParams) WithContext(ctx context.Context) *GetTasksCountParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get tasks count params
func (o *GetTasksCountParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get tasks count params
func (o *GetTasksCountParams) WithHTTPClient(client *http.Client) *GetTasksCountParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get tasks count params
func (o *GetTasksCountParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRepository adds the repository to the get tasks count params
func (o *GetTasksCountParams) WithRepository(repository *string) *GetTasksCountParams {
	o.SetRepository(repository)
	return o
}

// SetRepository adds the repository to the get tasks count params
func (o *GetTasksCountParams) SetRepository(repository *string) {
	o.Repository = repository
}

// WithSha adds the sha to the get tasks count params
func (o *GetTasksCountParams) WithSha(sha *string) *GetTasksCountParams {
	o.SetSha(sha)
	return o
}

// SetSha adds the sha to the get tasks count params
func (o *GetTasksCountParams) SetSha(sha *string) {
	o.Sha = sha
}

// WriteToRequest writes these params to a swagger request
func (o *GetTasksCountParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Repository != nil {

		o.HandleQueryParam(r, reg, "repository", o.Repository, nil)

	}

	if o.Sha != nil {

		o.HandleQueryParam(r, reg, "sha", o.Sha, nil)

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTasksCountParams) HandleQueryParam(r runtime.ClientRequest, reg strfmt.Registry, name string, param interface{}, formatter func(interface{}) string) error {
	if (reflect.TypeOf(param).Kind() == reflect.Ptr) && param == nil {
		return nil
	}

	if formatter == nil {
		formatter = func(i interface{}) string { return i.(string) }
	}

	return r.SetQueryParam(name, formatter(param))
}