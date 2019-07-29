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

// NewGetSubmissionIDTasksParams creates a new GetSubmissionIDTasksParams object
// with the default values initialized.
func NewGetSubmissionIDTasksParams() *GetSubmissionIDTasksParams {
	var (
		pageDefault    = int64(0)
		perPageDefault = int64(100)
	)
	return &GetSubmissionIDTasksParams{
		Page:    &pageDefault,
		PerPage: &perPageDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewGetSubmissionIDTasksParamsWithTimeout creates a new GetSubmissionIDTasksParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetSubmissionIDTasksParamsWithTimeout(timeout time.Duration) *GetSubmissionIDTasksParams {
	var (
		pageDefault    = int64(0)
		perPageDefault = int64(100)
	)
	return &GetSubmissionIDTasksParams{
		Page:    &pageDefault,
		PerPage: &perPageDefault,

		timeout: timeout,
	}
}

// NewGetSubmissionIDTasksParamsWithContext creates a new GetSubmissionIDTasksParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetSubmissionIDTasksParamsWithContext(ctx context.Context) *GetSubmissionIDTasksParams {
	var (
		pageDefault    = int64(0)
		perPageDefault = int64(100)
	)
	return &GetSubmissionIDTasksParams{
		Page:    &pageDefault,
		PerPage: &perPageDefault,

		Context: ctx,
	}
}

// NewGetSubmissionIDTasksParamsWithHTTPClient creates a new GetSubmissionIDTasksParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetSubmissionIDTasksParamsWithHTTPClient(client *http.Client) *GetSubmissionIDTasksParams {
	var (
		pageDefault    = int64(0)
		perPageDefault = int64(100)
	)
	return &GetSubmissionIDTasksParams{
		Page:       &pageDefault,
		PerPage:    &perPageDefault,
		HTTPClient: client,
	}
}

/*GetSubmissionIDTasksParams contains all the parameters to send to the API endpoint
for the get submission ID tasks operation typically these are written to a http.Request
*/
type GetSubmissionIDTasksParams struct {

	/*ID
	  The ID of the submission to retrieve

	*/
	ID int64
	/*Page
	  pagination control: what page to retrieve in the query.

	*/
	Page *int64
	/*PerPage
	  pagination control: how many items counts as a page.

	*/
	PerPage *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) WithTimeout(timeout time.Duration) *GetSubmissionIDTasksParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) WithContext(ctx context.Context) *GetSubmissionIDTasksParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) WithHTTPClient(client *http.Client) *GetSubmissionIDTasksParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) WithID(id int64) *GetSubmissionIDTasksParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) SetID(id int64) {
	o.ID = id
}

// WithPage adds the page to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) WithPage(page *int64) *GetSubmissionIDTasksParams {
	o.SetPage(page)
	return o
}

// SetPage adds the page to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) SetPage(page *int64) {
	o.Page = page
}

// WithPerPage adds the perPage to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) WithPerPage(perPage *int64) *GetSubmissionIDTasksParams {
	o.SetPerPage(perPage)
	return o
}

// SetPerPage adds the perPage to the get submission ID tasks params
func (o *GetSubmissionIDTasksParams) SetPerPage(perPage *int64) {
	o.PerPage = perPage
}

// WriteToRequest writes these params to a swagger request
func (o *GetSubmissionIDTasksParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if o.Page != nil {

		o.HandleQueryParam(r, reg, "page", o.Page, func(i interface{}) string { return swag.FormatInt64(i.(int64)) })

	}

	if o.PerPage != nil {

		o.HandleQueryParam(r, reg, "perPage", o.PerPage, func(i interface{}) string { return swag.FormatInt64(i.(int64)) })

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetSubmissionIDTasksParams) HandleQueryParam(r runtime.ClientRequest, reg strfmt.Registry, name string, param interface{}, formatter func(interface{}) string) error {
	if (reflect.TypeOf(param).Kind() == reflect.Ptr) && param == nil {
		return nil
	}

	if formatter == nil {
		formatter = func(i interface{}) string { return i.(string) }
	}

	return r.SetQueryParam(name, formatter(param))
}