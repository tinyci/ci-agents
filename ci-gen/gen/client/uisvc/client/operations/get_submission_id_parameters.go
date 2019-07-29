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
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetSubmissionIDParams creates a new GetSubmissionIDParams object
// with the default values initialized.
func NewGetSubmissionIDParams() *GetSubmissionIDParams {
	var ()
	return &GetSubmissionIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetSubmissionIDParamsWithTimeout creates a new GetSubmissionIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetSubmissionIDParamsWithTimeout(timeout time.Duration) *GetSubmissionIDParams {
	var ()
	return &GetSubmissionIDParams{

		timeout: timeout,
	}
}

// NewGetSubmissionIDParamsWithContext creates a new GetSubmissionIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetSubmissionIDParamsWithContext(ctx context.Context) *GetSubmissionIDParams {
	var ()
	return &GetSubmissionIDParams{

		Context: ctx,
	}
}

// NewGetSubmissionIDParamsWithHTTPClient creates a new GetSubmissionIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetSubmissionIDParamsWithHTTPClient(client *http.Client) *GetSubmissionIDParams {
	var ()
	return &GetSubmissionIDParams{
		HTTPClient: client,
	}
}

/*GetSubmissionIDParams contains all the parameters to send to the API endpoint
for the get submission ID operation typically these are written to a http.Request
*/
type GetSubmissionIDParams struct {

	/*ID
	  The ID of the submission to retrieve

	*/
	ID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get submission ID params
func (o *GetSubmissionIDParams) WithTimeout(timeout time.Duration) *GetSubmissionIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get submission ID params
func (o *GetSubmissionIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get submission ID params
func (o *GetSubmissionIDParams) WithContext(ctx context.Context) *GetSubmissionIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get submission ID params
func (o *GetSubmissionIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get submission ID params
func (o *GetSubmissionIDParams) WithHTTPClient(client *http.Client) *GetSubmissionIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get submission ID params
func (o *GetSubmissionIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get submission ID params
func (o *GetSubmissionIDParams) WithID(id int64) *GetSubmissionIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get submission ID params
func (o *GetSubmissionIDParams) SetID(id int64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetSubmissionIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetSubmissionIDParams) HandleQueryParam(r runtime.ClientRequest, reg strfmt.Registry, name string, param interface{}, formatter func(interface{}) string) error {
	if (reflect.TypeOf(param).Kind() == reflect.Ptr) && param == nil {
		return nil
	}

	if formatter == nil {
		formatter = func(i interface{}) string { return i.(string) }
	}

	return r.SetQueryParam(name, formatter(param))
}