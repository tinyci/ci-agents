package types

import (
	"fmt"
	"reflect"
	"time"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	yaml "gopkg.in/yaml.v2"
)

var (
	// ErrTaskParse is what is wrapped when a parse error occurs.
	ErrTaskParse = errors.New("parse error")
	// ErrTaskValidation is what is wrapped when errors in TaskSettings.Validate() occur.
	ErrTaskValidation = errors.New("validation error")
)

func mkMap(s *_struct.Struct) map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range s.Fields {
		m[k] = v.GetStringValue()
	}

	return m
}

func mkStruct(m map[string]interface{}) *_struct.Struct {
	s := &_struct.Struct{Fields: map[string]*_struct.Value{}}

	for k, v := range m {
		// FIXME make this typed
		s.Fields[k] = &_struct.Value{Kind: &_struct.Value_StringValue{StringValue: fmt.Sprintf("%v", v)}}
	}

	return s
}

// TaskSettings encompasses things that are a part of a task that are configurable by a user.
type TaskSettings struct {
	Mountpoint       string                  `yaml:"mountpoint"`
	Env              []string                `yaml:"env"`
	Dependencies     []string                `yaml:"dependencies"`
	WorkDir          string                  `yaml:"workdir"`
	Runs             map[string]*RunSettings `yaml:"runs"`
	DefaultTimeout   time.Duration           `yaml:"default_timeout"`
	DefaultQueue     string                  `yaml:"default_queue"`
	DefaultImage     string                  `yaml:"default_image"`
	Metadata         map[string]interface{}  `yaml:"metadata"`
	Config           RepoConfig              `yaml:"-"`
	DefaultResources Resources               `yaml:"default_resources"`
}

// NewTaskSettingsFromProto creates a task settings object from a proto representation.
func NewTaskSettingsFromProto(ts *types.TaskSettings) *TaskSettings {
	runs := map[string]*RunSettings{}

	for k, run := range ts.Runs {
		runs[k] = NewRunSettingsFromProto(run)
	}

	return &TaskSettings{
		Mountpoint:       ts.Mountpoint,
		Env:              ts.Env,
		WorkDir:          ts.Workdir,
		Runs:             runs,
		Dependencies:     ts.Dependencies,
		DefaultTimeout:   time.Duration(ts.DefaultTimeout),
		DefaultQueue:     ts.DefaultQueue,
		DefaultImage:     ts.DefaultImage,
		Metadata:         mkMap(ts.Metadata),
		DefaultResources: newResources(ts.Resources),
		Config:           NewRepoConfigFromProto(ts.Config),
	}
}

// ToProto makes a protobuf version of task settings.
func (t *TaskSettings) ToProto() *types.TaskSettings {
	runs := map[string]*types.RunSettings{}

	for name, run := range t.Runs {
		runs[name] = run.ToProto()
	}

	return &types.TaskSettings{
		Mountpoint:     t.Mountpoint,
		Env:            t.Env,
		Workdir:        t.WorkDir,
		Runs:           runs,
		Dependencies:   t.Dependencies,
		DefaultTimeout: t.DefaultTimeout.Nanoseconds(),
		DefaultQueue:   t.DefaultQueue,
		DefaultImage:   t.DefaultImage,
		Metadata:       mkStruct(t.Metadata),
		Resources:      t.DefaultResources.toProto(),
		Config:         t.Config.ToProto(),
	}
}

// NewTaskSettings creates a new task configuration from a byte buffer.
func NewTaskSettings(buf []byte, requireRuns bool, rc RepoConfig) (*TaskSettings, *errors.Error) {
	t := &TaskSettings{}

	if err := yaml.UnmarshalStrict(buf, t); err != nil {
		return nil, errors.New(err).Wrap(ErrTaskParse)
	}

	t.Config = rc

	if err := t.Validate(requireRuns); err != nil {
		return nil, errors.New(err).Wrap(ErrTaskValidation)
	}

	return t, nil
}

func (t *TaskSettings) overrideTimeout() {
	if t.Config.OverrideTimeout && t.Config.GlobalTimeout != 0 {
		for _, run := range t.Runs {
			run.Timeout = t.Config.GlobalTimeout
		}
	} else {
		for _, run := range t.Runs {
			if run.Timeout == 0 {
				if t.DefaultTimeout != 0 {
					run.Timeout = t.DefaultTimeout
				} else if t.Config.GlobalTimeout != 0 {
					run.Timeout = t.Config.GlobalTimeout
				}
			}
		}
	}
}

func (t *TaskSettings) overrideQueue() {
	if t.Config.OverrideQueue {
		for _, run := range t.Runs {
			run.Queue = t.Config.Queue
		}
	} else if t.DefaultQueue != "" {
		for _, run := range t.Runs {
			if run.Queue == "" {
				run.Queue = t.DefaultQueue
			}
		}
	} else {
		for _, run := range t.Runs {
			if run.Queue == "" {
				if t.Config.Queue != "" {
					run.Queue = t.Config.Queue
				} else {
					run.Queue = "default"
				}
			}
		}
	}
}

func (t *TaskSettings) handleOverrides() {
	if t.DefaultImage == "" && t.Config.DefaultImage != "" {
		t.DefaultImage = t.Config.DefaultImage
	}

	for name, run := range t.Runs {
		if run.Image == "" && t.DefaultImage != "" {
			run.Image = t.DefaultImage
		}

		// propagate repoconfig metadata into the runs
		for key, value := range t.Config.Metadata {
			if run.Metadata == nil {
				run.Metadata = map[string]interface{}{}
			}

			if _, ok := run.Metadata[key]; !ok || t.Config.OverrideMetadata {
				run.Metadata[key] = value
			}
		}
		run.Name = name
	}

	t.overrideTimeout()
	t.overrideQueue()

	if t.WorkDir == "" && t.Config.WorkDir != "" {
		t.WorkDir = t.Config.WorkDir
	}

	if reflect.DeepEqual(t.DefaultResources, Resources{}) && !reflect.DeepEqual(t.Config.DefaultResources, Resources{}) {
		t.Config.DefaultResources.copy(&t.DefaultResources)
	}
}

// Validate validates the task settings.
func (t *TaskSettings) Validate(requireRuns bool) *errors.Error {
	t.handleOverrides()

	if !t.Config.AllowPrivileged {
		for _, run := range t.Runs {
			if run.Privileged {
				return errors.Errorf("Run %q cannot launch because it wants a privileged container and they are denied", run.Name)
			}
		}
	}

	if len(t.Runs) != 0 {
		if t.Mountpoint == "" {
			return errors.New("no mountpoint")
		}

		if t.WorkDir == "" {
			t.WorkDir = t.Mountpoint
		}

		for _, run := range t.Runs {
			if err := run.Validate(t); err != nil {
				return err
			}
		}
	} else if requireRuns {
		return errors.New("Runs are required to proceed further with this task")
	} else if len(t.Dependencies) == 0 {
		return errors.New("no runs in task and no dependencies")
	}

	return nil
}

// RunSettings encompasses things that are a part of a run that are
// configurable by a user.
type RunSettings struct {
	Privileged bool                   `yaml:"privileged"`
	Command    []string               `yaml:"command"`
	Image      string                 `yaml:"image"`
	Queue      string                 `yaml:"queue"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Name       string                 `yaml:"-"`
	Timeout    time.Duration          `yaml:"timeout"`
	Resources  Resources              `yaml:"resources"`
	Env        []string               `yaml:"env"`
}

// Resources communicates what resources should be available to the runner.
// This can be overridden by the tinyci master configuration.
//
// The values here and their measurements are interpreted by the runner. We
// do not normalize between runners.
type Resources struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
	Disk   string `yaml:"disk"`
	IOPS   string `yaml:"iops"`
}

func (r Resources) toProto() *types.Resources {
	return &types.Resources{
		Cpu:    r.CPU,
		Memory: r.Memory,
		Disk:   r.Disk,
		Iops:   r.IOPS,
	}
}

func (r Resources) copy(r2 *Resources) {
	r2.CPU = r.CPU
	r2.Memory = r.Memory
	r2.Disk = r.Disk
	r2.IOPS = r.IOPS
}

func newResources(rs *types.Resources) Resources {
	return Resources{
		CPU:    rs.Cpu,
		Memory: rs.Memory,
		Disk:   rs.Disk,
		IOPS:   rs.Iops,
	}
}

// NewRunSettingsFromProto creates a runsettings from a proto representation.
func NewRunSettingsFromProto(rs *types.RunSettings) *RunSettings {
	return &RunSettings{
		Privileged: rs.Privileged,
		Command:    rs.Command,
		Image:      rs.Image,
		Queue:      rs.Queue,
		Metadata:   mkMap(rs.Metadata),
		Name:       rs.Name,
		Timeout:    time.Duration(rs.Timeout),
		Resources:  newResources(rs.Resources),
		Env:        rs.Env,
	}
}

// ToProto converts the run settings to the protobuf representation.
func (rs *RunSettings) ToProto() *types.RunSettings {
	return &types.RunSettings{
		Privileged: rs.Privileged,
		Command:    rs.Command,
		Image:      rs.Image,
		Queue:      rs.Queue,
		Metadata:   mkStruct(rs.Metadata),
		Name:       rs.Name,
		Timeout:    rs.Timeout.Nanoseconds(),
		Resources:  rs.Resources.toProto(),
		Env:        rs.Env,
	}
}

// Validate validates the run settings, returning errors on any found.
func (rs *RunSettings) Validate(t *TaskSettings) *errors.Error {
	if len(rs.Command) == 0 {
		return errors.New("command was empty")
	}

	if rs.Image == "" {
		return errors.New("image was empty")
	}

	if rs.Queue == "" {
		return errors.New("queue name was empty")
	}

	if reflect.DeepEqual(rs.Resources, Resources{}) {
		t.DefaultResources.copy(&rs.Resources)
	}

	return nil
}

// RepoConfig is the global configuration for the repository. It allows setting
// of global attributes as well as overrides and defaults for certain
// task-related items. It is typically named `tinyci.yml`.
type RepoConfig struct {
	AllowPrivileged  bool                   `yaml:"allow_privileged"`
	WorkDir          string                 `yaml:"workdir"`
	Queue            string                 `yaml:"queue"`
	OverrideQueue    bool                   `yaml:"override_queue"`
	GlobalTimeout    time.Duration          `yaml:"global_timeout"` // run timeout. if unset, or 0, no timeout.
	OverrideTimeout  bool                   `yaml:"override_timeout"`
	IgnoreDirs       []string               `yaml:"ignore_directories"`
	Metadata         map[string]interface{} `yaml:"metadata"`
	OverrideMetadata bool                   `yaml:"override_metadata"`
	DefaultImage     string                 `yaml:"default_image"`
	DefaultResources Resources              `yaml:"default_resources"`
	//OptimizeDiff  bool   `yaml:"optimize_diff"` // diff dir selection -- FIXME defaulted to on for now, will add this logic later
}

// NewRepoConfig creates a new repo config from a byte buffer.
func NewRepoConfig(buf []byte) (RepoConfig, *errors.Error) {
	r := RepoConfig{}

	if err := yaml.UnmarshalStrict(buf, &r); err != nil {
		return r, errors.New(err)
	}

	r.handleOverrides()

	return r, r.Validate()
}

// NewRepoConfigFromProto creates a runsettings from a proto representation.
func NewRepoConfigFromProto(rs *types.RepoConfig) RepoConfig {
	metadata := map[string]interface{}{}

	for key, value := range rs.Metadata {
		metadata[key] = value
	}

	return RepoConfig{
		AllowPrivileged:  rs.AllowPrivileged,
		WorkDir:          rs.Workdir,
		Queue:            rs.Queue,
		OverrideQueue:    rs.OverrideQueue,
		GlobalTimeout:    time.Duration(rs.GlobalTimeout),
		OverrideTimeout:  rs.OverrideTimeout,
		IgnoreDirs:       rs.IgnoreDirectories,
		Metadata:         metadata,
		OverrideMetadata: rs.OverrideMetadata,
		DefaultImage:     rs.DefaultImage,
		DefaultResources: newResources(rs.DefaultResources),
	}
}

// ToProto converts the run settings to the protobuf representation.
func (r *RepoConfig) ToProto() *types.RepoConfig {
	metadata := map[string]string{}

	for key, value := range r.Metadata {
		s, ok := value.(string)
		if ok {
			metadata[key] = s
		}
	}

	return &types.RepoConfig{
		AllowPrivileged:   r.AllowPrivileged,
		Workdir:           r.WorkDir,
		Queue:             r.Queue,
		OverrideQueue:     r.OverrideQueue,
		GlobalTimeout:     int64(r.GlobalTimeout),
		OverrideTimeout:   r.OverrideTimeout,
		IgnoreDirectories: r.IgnoreDirs,
		Metadata:          metadata,
		OverrideMetadata:  r.OverrideMetadata,
		DefaultImage:      r.DefaultImage,
		DefaultResources:  r.DefaultResources.toProto(),
	}
}

func (r *RepoConfig) handleOverrides() {
	if r.Queue == "" {
		r.Queue = "default"
	}
}

// Validate returns any error if there are validation errors in the repo config.
func (r *RepoConfig) Validate() *errors.Error {
	// it doesn't do much right now..
	if r.Queue == "" {
		return errors.New("queue was empty")
	}

	return nil
}
