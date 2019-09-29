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

// TaskList task list
// swagger:model taskList
type TaskList []*TaskListItems0

// Validate validates this task list
func (m TaskList) Validate(formats strfmt.Registry) error {
	var res []error

	for i := 0; i < len(m); i++ {
		if swag.IsZero(m[i]) { // not required
			continue
		}

		if m[i] != nil {
			if err := m[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName(strconv.Itoa(i))
				}
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// TaskListItems0 task list items0
// swagger:model TaskListItems0
type TaskListItems0 struct {

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
	Parent *TaskListItems0Parent `json:"parent,omitempty"`

	// path
	Path string `json:"path,omitempty"`

	// ref
	Ref *TaskListItems0Ref `json:"ref,omitempty"`

	// runs
	Runs int64 `json:"runs,omitempty"`

	// settings
	Settings *TaskListItems0Settings `json:"settings,omitempty"`

	// started at
	// Format: date-time
	StartedAt *strfmt.DateTime `json:"started_at,omitempty"`

	// status
	Status *bool `json:"status,omitempty"`
}

// Validate validates this task list items0
func (m *TaskListItems0) Validate(formats strfmt.Registry) error {
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

func (m *TaskListItems0) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *TaskListItems0) validateFinishedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.FinishedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("finished_at", "body", "date-time", m.FinishedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *TaskListItems0) validateParent(formats strfmt.Registry) error {

	if swag.IsZero(m.Parent) { // not required
		return nil
	}

	if m.Parent != nil {
		if err := m.Parent.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("parent")
			}
			return err
		}
	}

	return nil
}

func (m *TaskListItems0) validateRef(formats strfmt.Registry) error {

	if swag.IsZero(m.Ref) { // not required
		return nil
	}

	if m.Ref != nil {
		if err := m.Ref.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ref")
			}
			return err
		}
	}

	return nil
}

func (m *TaskListItems0) validateSettings(formats strfmt.Registry) error {

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

func (m *TaskListItems0) validateStartedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.StartedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("started_at", "body", "date-time", m.StartedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TaskListItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskListItems0) UnmarshalBinary(b []byte) error {
	var res TaskListItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// TaskListItems0Parent task list items0 parent
// swagger:model TaskListItems0Parent
type TaskListItems0Parent struct {

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

// Validate validates this task list items0 parent
func (m *TaskListItems0Parent) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TaskListItems0Parent) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskListItems0Parent) UnmarshalBinary(b []byte) error {
	var res TaskListItems0Parent
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// TaskListItems0Ref task list items0 ref
// swagger:model TaskListItems0Ref
type TaskListItems0Ref struct {

	// id
	ID int64 `json:"id,omitempty"`

	// ref name
	RefName string `json:"ref_name,omitempty"`

	// repository
	Repository *TaskListItems0RefRepository `json:"repository,omitempty"`

	// sha
	Sha string `json:"sha,omitempty"`
}

// Validate validates this task list items0 ref
func (m *TaskListItems0Ref) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRepository(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TaskListItems0Ref) validateRepository(formats strfmt.Registry) error {

	if swag.IsZero(m.Repository) { // not required
		return nil
	}

	if m.Repository != nil {
		if err := m.Repository.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ref" + "." + "repository")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TaskListItems0Ref) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskListItems0Ref) UnmarshalBinary(b []byte) error {
	var res TaskListItems0Ref
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// TaskListItems0RefRepository task list items0 ref repository
// swagger:model TaskListItems0RefRepository
type TaskListItems0RefRepository struct {

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

// Validate validates this task list items0 ref repository
func (m *TaskListItems0RefRepository) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TaskListItems0RefRepository) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskListItems0RefRepository) UnmarshalBinary(b []byte) error {
	var res TaskListItems0RefRepository
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// TaskListItems0Settings task list items0 settings
// swagger:model TaskListItems0Settings
type TaskListItems0Settings struct {

	// config
	Config *TaskListItems0SettingsConfig `json:"config,omitempty"`

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
	Runs map[string]TaskListItems0SettingsRunsAnon `json:"runs,omitempty"`

	// workdir
	Workdir string `json:"workdir,omitempty"`
}

// Validate validates this task list items0 settings
func (m *TaskListItems0Settings) Validate(formats strfmt.Registry) error {
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

func (m *TaskListItems0Settings) validateConfig(formats strfmt.Registry) error {

	if swag.IsZero(m.Config) { // not required
		return nil
	}

	if m.Config != nil {
		if err := m.Config.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("settings" + "." + "config")
			}
			return err
		}
	}

	return nil
}

func (m *TaskListItems0Settings) validateRuns(formats strfmt.Registry) error {

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
func (m *TaskListItems0Settings) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskListItems0Settings) UnmarshalBinary(b []byte) error {
	var res TaskListItems0Settings
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// TaskListItems0SettingsConfig task list items0 settings config
// swagger:model TaskListItems0SettingsConfig
type TaskListItems0SettingsConfig struct {

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

// Validate validates this task list items0 settings config
func (m *TaskListItems0SettingsConfig) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TaskListItems0SettingsConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskListItems0SettingsConfig) UnmarshalBinary(b []byte) error {
	var res TaskListItems0SettingsConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// TaskListItems0SettingsRunsAnon task list items0 settings runs anon
// swagger:model TaskListItems0SettingsRunsAnon
type TaskListItems0SettingsRunsAnon struct {

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

// Validate validates this task list items0 settings runs anon
func (m *TaskListItems0SettingsRunsAnon) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TaskListItems0SettingsRunsAnon) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskListItems0SettingsRunsAnon) UnmarshalBinary(b []byte) error {
	var res TaskListItems0SettingsRunsAnon
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
