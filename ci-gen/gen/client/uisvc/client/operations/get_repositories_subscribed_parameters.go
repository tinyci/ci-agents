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

// NewGetRepositoriesSubscribedParams creates a new GetRepositoriesSubscribedParams object
// with the default values initialized.
func NewGetRepositoriesSubscribedParams() *GetRepositoriesSubscribedParams {
	var ()
	return &GetRepositoriesSubscribedParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetRepositoriesSubscribedParamsWithTimeout creates a new GetRepositoriesSubscribedParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetRepositoriesSubscribedParamsWithTimeout(timeout time.Duration) *GetRepositoriesSubscribedParams {
	var ()
	return &GetRepositoriesSubscribedParams{

		timeout: timeout,
	}
}

// NewGetRepositoriesSubscribedParamsWithContext creates a new GetRepositoriesSubscribedParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetRepositoriesSubscribedParamsWithContext(ctx context.Context) *GetRepositoriesSubscribedParams {
	var ()
	return &GetRepositoriesSubscribedParams{

		Context: ctx,
	}
}

// NewGetRepositoriesSubscribedParamsWithHTTPClient creates a new GetRepositoriesSubscribedParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetRepositoriesSubscribedParamsWithHTTPClient(client *http.Client) *GetRepositoriesSubscribedParams {
	var ()
	return &GetRepositoriesSubscribedParams{
		HTTPClient: client,
	}
}

/*GetRepositoriesSubscribedParams contains all the parameters to send to the API endpoint
for the get repositories subscribed operation typically these are written to a http.Request
*/
type GetRepositoriesSubscribedParams struct {

	/*Search
	  search string by which to filter results

	*/
	Search *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get repositories subscribed params
func (o *GetRepositoriesSubscribedParams) WithTimeout(timeout time.Duration) *GetRepositoriesSubscribedParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get repositories subscribed params
func (o *GetRepositoriesSubscribedParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get repositories subscribed params
func (o *GetRepositoriesSubscribedParams) WithContext(ctx context.Context) *GetRepositoriesSubscribedParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get repositories subscribed params
func (o *GetRepositoriesSubscribedParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get repositories subscribed params
func (o *GetRepositoriesSubscribedParams) WithHTTPClient(client *http.Client) *GetRepositoriesSubscribedParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get repositories subscribed params
func (o *GetRepositoriesSubscribedParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithSearch adds the search to the get repositories subscribed params
func (o *GetRepositoriesSubscribedParams) WithSearch(search *string) *GetRepositoriesSubscribedParams {
	o.SetSearch(search)
	return o
}

// SetSearch adds the search to the get repositories subscribed params
func (o *GetRepositoriesSubscribedParams) SetSearch(search *string) {
	o.Search = search
}

// WriteToRequest writes these params to a swagger request
func (o *GetRepositoriesSubscribedParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Search != nil {

		o.HandleQueryParam(r, reg, "search", o.Search, nil)

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetRepositoriesSubscribedParams) HandleQueryParam(r runtime.ClientRequest, reg strfmt.Registry, name string, param interface{}, formatter func(interface{}) string) error {
	if (reflect.TypeOf(param).Kind() == reflect.Ptr) && param == nil {
		return nil
	}

	if formatter == nil {
		formatter = func(i interface{}) string { return i.(string) }
	}

	return r.SetQueryParam(name, formatter(param))
}