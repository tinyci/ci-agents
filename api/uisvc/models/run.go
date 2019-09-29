// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Run run
// swagger:model run
type Run struct {

	// created at
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"created_at,omitempty"`

	// finished at
	// Format: date-time
	FinishedAt *strfmt.DateTime `json:"finished_at,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// ran on
	RanOn *string `json:"ran_on,omitempty"`

	// settings
	Settings *RunSettings `json:"settings,omitempty"`

	// started at
	// Format: date-time
	StartedAt *strfmt.DateTime `json:"started_at,omitempty"`

	// status
	Status *bool `json:"status,omitempty"`

	// task
	Task *RunTask `json:"task,omitempty"`
}

// Validate validates this run
func (m *Run) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFinishedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSettings(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTask(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Run) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Run) validateFinishedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.FinishedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("finished_at", "body", "date-time", m.FinishedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Run) validateSettings(formats strfmt.Registry) error {

	if swag.IsZero(m.Settings) { // not required
		return nil
	}

	if m.Settings != nil {
		if err := m.Settings.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("settings")
			}
			return err
		}
	}

	return nil
}

func (m *Run) validateStartedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.StartedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("started_at", "body", "date-time", m.StartedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Run) validateTask(formats strfmt.Registry) error {

	if swag.IsZero(m.Task) { // not required
		return nil
	}

	if m.Task != nil {
		if err := m.Task.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Run) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Run) UnmarshalBinary(b []byte) error {
	var res Run
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// RunSettings run settings
// swagger:model RunSettings
type RunSettings struct {

	// command
	Command []string `json:"command"`

	// image
	Image string `json:"image,omitempty"`

	// metadata
	Metadata interface{} `json:"metadata,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// queue
	Queue string `json:"queue,omitempty"`

	// timeout
	Timeout int64 `json:"timeout,omitempty"`
}

// Validate validates this run settings
func (m *RunSettings) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RunSettings) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RunSettings) UnmarshalBinary(b []byte) error {
	var res RunSettings
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// RunTask run task
// swagger:model RunTask
type RunTask struct {

	// base sha
	BaseSha string `json:"base_sha,omitempty"`

	// canceled
	Canceled bool `json:"canceled,omitempty"`

	// created at
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"created_at,omitempty"`

	// finished at
	// Format: date-time
	FinishedAt *strfmt.DateTime `json:"finished_at,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// parent
	Parent *RunTaskParent `json:"parent,omitempty"`

	// path
	Path string `json:"path,omitempty"`

	// ref
	Ref *RunTaskRef `json:"ref,omitempty"`

	// runs
	Runs int64 `json:"runs,omitempty"`

	// settings
	Settings *RunTaskSettings `json:"settings,omitempty"`

	// started at
	// Format: date-time
	StartedAt *strfmt.DateTime `json:"started_at,omitempty"`

	// status
	Status *bool `json:"status,omitempty"`
}

// Validate validates this run task
func (m *RunTask) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFinishedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParent(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRef(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSettings(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RunTask) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("task"+"."+"created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *RunTask) validateFinishedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.FinishedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("task"+"."+"finished_at", "body", "date-time", m.FinishedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *RunTask) validateParent(formats strfmt.Registry) error {

	if swag.IsZero(m.Parent) { // not required
		return nil
	}

	if m.Parent != nil {
		if err := m.Parent.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task" + "." + "parent")
			}
			return err
		}
	}

	return nil
}

func (m *RunTask) validateRef(formats strfmt.Registry) error {

	if swag.IsZero(m.Ref) { // not required
		return nil
	}

	if m.Ref != nil {
		if err := m.Ref.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task" + "." + "ref")
			}
			return err
		}
	}

	return nil
}

func (m *RunTask) validateSettings(formats strfmt.Registry) error {

	if swag.IsZero(m.Settings) { // not required
		return nil
	}

	if m.Settings != nil {
		if err := m.Settings.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task" + "." + "settings")
			}
			return err
		}
	}

	return nil
}

func (m *RunTask) validateStartedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.StartedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("task"+"."+"started_at", "body", "date-time", m.StartedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RunTask) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RunTask) UnmarshalBinary(b []byte) error {
	var res RunTask
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// RunTaskParent run task parent
// swagger:model RunTaskParent
type RunTaskParent struct {

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

// Validate validates this run task parent
func (m *RunTaskParent) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RunTaskParent) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RunTaskParent) UnmarshalBinary(b []byte) error {
	var res RunTaskParent
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// RunTaskRef run task ref
// swagger:model RunTaskRef
type RunTaskRef struct {

	// id
	ID int64 `json:"id,omitempty"`

	// ref name
	RefName string `json:"ref_name,omitempty"`

	// repository
	Repository *RunTaskRefRepository `json:"repository,omitempty"`

	// sha
	Sha string `json:"sha,omitempty"`
}

// Validate validates this run task ref
func (m *RunTaskRef) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRepository(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RunTaskRef) validateRepository(formats strfmt.Registry) error {

	if swag.IsZero(m.Repository) { // not required
		return nil
	}

	if m.Repository != nil {
		if err := m.Repository.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task" + "." + "ref" + "." + "repository")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RunTaskRef) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RunTaskRef) UnmarshalBinary(b []byte) error {
	var res RunTaskRef
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// RunTaskRefRepository run task ref repository
// swagger:model RunTaskRefRepository
type RunTaskRefRepository struct {

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

// Validate validates this run task ref repository
func (m *RunTaskRefRepository) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RunTaskRefRepository) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RunTaskRefRepository) UnmarshalBinary(b []byte) error {
	var res RunTaskRefRepository
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// RunTaskSettings run task settings
// swagger:model RunTaskSettings
type RunTaskSettings struct {

	// config
	Config *RunTaskSettingsConfig `json:"config,omitempty"`

	// default image
	DefaultImage string `json:"default_image,omitempty"`

	// default queue
	DefaultQueue string `json:"default_queue,omitempty"`

	// the default timeout; in nanoseconds
	DefaultTimeout int64 `json:"default_timeout,omitempty"`

	// env
	Env []string `json:"env"`

	// metadata
	Metadata interface{} `json:"metadata,omitempty"`

	// mountpoint
	Mountpoint string `json:"mountpoint,omitempty"`

	// runs
	Runs map[string]RunTaskSettingsRunsAnon `json:"runs,omitempty"`

	// workdir
	Workdir string `json:"workdir,omitempty"`
}

// Validate validates this run task settings
func (m *RunTaskSettings) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConfig(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRuns(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RunTaskSettings) validateConfig(formats strfmt.Registry) error {

	if swag.IsZero(m.Config) { // not required
		return nil
	}

	if m.Config != nil {
		if err := m.Config.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task" + "." + "settings" + "." + "config")
			}
			return err
		}
	}

	return nil
}

func (m *RunTaskSettings) validateRuns(formats strfmt.Registry) error {

	if swag.IsZero(m.Runs) { // not required
		return nil
	}

	for k := range m.Runs {

		if swag.IsZero(m.Runs[k]) { // not required
			continue
		}
		if val, ok := m.Runs[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *RunTaskSettings) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RunTaskSettings) UnmarshalBinary(b []byte) error {
	var res RunTaskSettings
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// RunTaskSettingsConfig run task settings config
// swagger:model RunTaskSettingsConfig
type RunTaskSettingsConfig struct {

	// global timeout
	GlobalTimeout int64 `json:"global_timeout,omitempty"`

	// override queue
	OverrideQueue bool `json:"override_queue,omitempty"`

	// override timeout
	OverrideTimeout bool `json:"override_timeout,omitempty"`

	// queue
	Queue string `json:"queue,omitempty"`

	// workdir
	Workdir string `json:"workdir,omitempty"`
}

// Validate validates this run task settings config
func (m *RunTaskSettingsConfig) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RunTaskSettingsConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RunTaskSettingsConfig) UnmarshalBinary(b []byte) error {
	var res RunTaskSettingsConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// RunTaskSettingsRunsAnon run task settings runs anon
// swagger:model RunTaskSettingsRunsAnon
type RunTaskSettingsRunsAnon struct {

	// command
	Command []string `json:"command"`

	// image
	Image string `json:"image,omitempty"`

	// metadata
	Metadata interface{} `json:"metadata,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// queue
	Queue string `json:"queue,omitempty"`

	// timeout
	Timeout int64 `json:"timeout,omitempty"`
}

// Validate validates this run task settings runs anon
func (m *RunTaskSettingsRunsAnon) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RunTaskSettingsRunsAnon) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RunTaskSettingsRunsAnon) UnmarshalBinary(b []byte) error {
	var res RunTaskSettingsRunsAnon
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
