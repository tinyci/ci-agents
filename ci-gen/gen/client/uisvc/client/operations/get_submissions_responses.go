// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/tinyci/ci-agents/ci-gen/gen/client/uisvc/models"
)

// GetSubmissionsReader is a Reader for the GetSubmissions structure.
type GetSubmissionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSubmissionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetSubmissionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 500:
		result := NewGetSubmissionsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetSubmissionsOK creates a GetSubmissionsOK with default headers values
func NewGetSubmissionsOK() *GetSubmissionsOK {
	return &GetSubmissionsOK{}
}

/*GetSubmissionsOK handles this case with default header values.

OK
*/
type GetSubmissionsOK struct {
	Payload models.ModelSubmissionList
}

func (o *GetSubmissionsOK) Error() string {
	return fmt.Sprintf("[GET /submissions][%d] getSubmissionsOK  %+v", 200, o.Payload) // #nosec
}

func (o *GetSubmissionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSubmissionsInternalServerError creates a GetSubmissionsInternalServerError with default headers values
func NewGetSubmissionsInternalServerError() *GetSubmissionsInternalServerError {
	return &GetSubmissionsInternalServerError{}
}

/*GetSubmissionsInternalServerError handles this case with default header values.

Internal Server Error
*/
type GetSubmissionsInternalServerError struct {
	Payload *models.Error
}

func (o *GetSubmissionsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /submissions][%d] getSubmissionsInternalServerError  %+v", 500, o.Payload) // #nosec
}

func (o *GetSubmissionsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}