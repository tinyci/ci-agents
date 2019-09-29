// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ModelSubmission model submission
// swagger:model modelSubmission
type ModelSubmission struct {

	// base ref
	BaseRef *ModelSubmissionBaseRef `json:"base_ref,omitempty"`

	// canceled
	Canceled bool `json:"canceled,omitempty"`

	// created at
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"created_at,omitempty"`

	// finished at
	// Format: date-time
	FinishedAt *strfmt.DateTime `json:"finished_at,omitempty"`

	// head ref
	HeadRef *ModelSubmissionHeadRef `json:"head_ref,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// started at
	// Format: date-time
	StartedAt *strfmt.DateTime `json:"started_at,omitempty"`

	// status
	Status *bool `json:"status,omitempty"`

	// tasks count
	TasksCount int64 `json:"tasks_count,omitempty"`

	// ticket id
	TicketID int64 `json:"ticket_id,omitempty"`

	// user
	User *ModelSubmissionUser `json:"user,omitempty"`
}

// Validate validates this model submission
func (m *ModelSubmission) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBaseRef(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFinishedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHeadRef(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUser(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModelSubmission) validateBaseRef(formats strfmt.Registry) error {

	if swag.IsZero(m.BaseRef) { // not required
		return nil
	}

	if m.BaseRef != nil {
		if err := m.BaseRef.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("base_ref")
			}
			return err
		}
	}

	return nil
}

func (m *ModelSubmission) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ModelSubmission) validateFinishedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.FinishedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("finished_at", "body", "date-time", m.FinishedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ModelSubmission) validateHeadRef(formats strfmt.Registry) error {

	if swag.IsZero(m.HeadRef) { // not required
		return nil
	}

	if m.HeadRef != nil {
		if err := m.HeadRef.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("head_ref")
			}
			return err
		}
	}

	return nil
}

func (m *ModelSubmission) validateStartedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.StartedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("started_at", "body", "date-time", m.StartedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ModelSubmission) validateUser(formats strfmt.Registry) error {

	if swag.IsZero(m.User) { // not required
		return nil
	}

	if m.User != nil {
		if err := m.User.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("user")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModelSubmission) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelSubmission) UnmarshalBinary(b []byte) error {
	var res ModelSubmission
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ModelSubmissionBaseRef model submission base ref
// swagger:model ModelSubmissionBaseRef
type ModelSubmissionBaseRef struct {

	// id
	ID int64 `json:"id,omitempty"`

	// ref name
	RefName string `json:"ref_name,omitempty"`

	// repository
	Repository *ModelSubmissionBaseRefRepository `json:"repository,omitempty"`

	// sha
	Sha string `json:"sha,omitempty"`
}

// Validate validates this model submission base ref
func (m *ModelSubmissionBaseRef) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRepository(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModelSubmissionBaseRef) validateRepository(formats strfmt.Registry) error {

	if swag.IsZero(m.Repository) { // not required
		return nil
	}

	if m.Repository != nil {
		if err := m.Repository.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("base_ref" + "." + "repository")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModelSubmissionBaseRef) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelSubmissionBaseRef) UnmarshalBinary(b []byte) error {
	var res ModelSubmissionBaseRef
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ModelSubmissionBaseRefRepository model submission base ref repository
// swagger:model ModelSubmissionBaseRefRepository
type ModelSubmissionBaseRefRepository struct {

	// auto created
	AutoCreated bool `json:"auto_created,omitempty"`

	// disabled
	Disabled bool `json:"disabled,omitempty"`

	// github
	Github interface{} `json:"github,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// private
	Private bool `json:"private,omitempty"`
}

// Validate validates this model submission base ref repository
func (m *ModelSubmissionBaseRefRepository) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModelSubmissionBaseRefRepository) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelSubmissionBaseRefRepository) UnmarshalBinary(b []byte) error {
	var res ModelSubmissionBaseRefRepository
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ModelSubmissionHeadRef model submission head ref
// swagger:model ModelSubmissionHeadRef
type ModelSubmissionHeadRef struct {

	// id
	ID int64 `json:"id,omitempty"`

	// ref name
	RefName string `json:"ref_name,omitempty"`

	// repository
	Repository *ModelSubmissionHeadRefRepository `json:"repository,omitempty"`

	// sha
	Sha string `json:"sha,omitempty"`
}

// Validate validates this model submission head ref
func (m *ModelSubmissionHeadRef) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRepository(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModelSubmissionHeadRef) validateRepository(formats strfmt.Registry) error {

	if swag.IsZero(m.Repository) { // not required
		return nil
	}

	if m.Repository != nil {
		if err := m.Repository.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("head_ref" + "." + "repository")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModelSubmissionHeadRef) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelSubmissionHeadRef) UnmarshalBinary(b []byte) error {
	var res ModelSubmissionHeadRef
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ModelSubmissionHeadRefRepository model submission head ref repository
// swagger:model ModelSubmissionHeadRefRepository
type ModelSubmissionHeadRefRepository struct {

	// auto created
	AutoCreated bool `json:"auto_created,omitempty"`

	// disabled
	Disabled bool `json:"disabled,omitempty"`

	// github
	Github interface{} `json:"github,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// private
	Private bool `json:"private,omitempty"`
}

// Validate validates this model submission head ref repository
func (m *ModelSubmissionHeadRefRepository) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModelSubmissionHeadRefRepository) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelSubmissionHeadRefRepository) UnmarshalBinary(b []byte) error {
	var res ModelSubmissionHeadRefRepository
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ModelSubmissionUser model submission user
// swagger:model ModelSubmissionUser
type ModelSubmissionUser struct {

	// errors
	Errors []*ModelSubmissionUserErrorsItems0 `json:"errors"`

	// id
	ID int64 `json:"id,omitempty"`

	// last scanned repos
	// Format: date-time
	LastScannedRepos *strfmt.DateTime `json:"last_scanned_repos,omitempty"`

	// token
	Token interface{} `json:"token,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this model submission user
func (m *ModelSubmissionUser) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateErrors(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastScannedRepos(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModelSubmissionUser) validateErrors(formats strfmt.Registry) error {

	if swag.IsZero(m.Errors) { // not required
		return nil
	}

	for i := 0; i < len(m.Errors); i++ {
		if swag.IsZero(m.Errors[i]) { // not required
			continue
		}

		if m.Errors[i] != nil {
			if err := m.Errors[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("user" + "." + "errors" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ModelSubmissionUser) validateLastScannedRepos(formats strfmt.Registry) error {

	if swag.IsZero(m.LastScannedRepos) { // not required
		return nil
	}

	if err := validate.FormatOf("user"+"."+"last_scanned_repos", "body", "date-time", m.LastScannedRepos.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModelSubmissionUser) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelSubmissionUser) UnmarshalBinary(b []byte) error {
	var res ModelSubmissionUser
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ModelSubmissionUserErrorsItems0 model submission user errors items0
// swagger:model ModelSubmissionUserErrorsItems0
type ModelSubmissionUserErrorsItems0 struct {

	// error
	Error string `json:"error,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`
}

// Validate validates this model submission user errors items0
func (m *ModelSubmissionUserErrorsItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModelSubmissionUserErrorsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelSubmissionUserErrorsItems0) UnmarshalBinary(b []byte) error {
	var res ModelSubmissionUserErrorsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
