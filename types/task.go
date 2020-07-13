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
	Config           *RepoConfig             `yaml:"-"`
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
	}
}

// NewTaskSettings creates a new task configuration from a byte buffer.
func NewTaskSettings(buf []byte, requireRuns bool, rc *RepoConfig) (*TaskSettings, *errors.Error) {
	if rc == nil {
		rc = &RepoConfig{}
	}

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
	if t.Config != nil && t.Config.OverrideTimeout && t.Config.GlobalTimeout != 0 {
		for _, run := range t.Runs {
			run.Timeout = t.Config.GlobalTimeout
		}
	} else {
		for _, run := range t.Runs {
			if run.Timeout == 0 {
				if t.DefaultTimeout != 0 {
					run.Timeout = t.DefaultTimeout
				} else if t.Config != nil && t.Config.GlobalTimeout != 0 {
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
		t.DefaultResources = t.Config.DefaultResources
	}
}

// Validate validates the task settings.
func (t *TaskSettings) Validate(requireRuns bool) *errors.Error {
	if t.Config == nil {
		t.Config = &RepoConfig{}
	}

	t.handleOverrides()

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
	Command   []string               `yaml:"command"`
	Image     string                 `yaml:"image"`
	Queue     string                 `yaml:"queue"`
	Metadata  map[string]interface{} `yaml:"metadata"`
	Name      string                 `yaml:"-"`
	Timeout   time.Duration          `yaml:"timeout"`
	Resources Resources              `yaml:"resources"`
}

// Resources communicates what resources should be available to the runner.
// This can be overridden by the tinyci master configuration.
//
// The values here and their measurements are interpreted by the runner. We
// do not normalize between runners.
type Resources struct {
	CPU    uint32 `yaml:"cpu"`
	Memory uint32 `yaml:"memory"`
	Disk   uint32 `yaml:"disk"`
	IOPS   uint32 `yaml:"iops"`
}

func (r Resources) toProto() *types.Resources {
	return &types.Resources{
		Cpu:    r.CPU,
		Memory: r.Memory,
		Disk:   r.Disk,
		Iops:   r.IOPS,
	}
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
		Command:   rs.Command,
		Image:     rs.Image,
		Queue:     rs.Queue,
		Metadata:  mkMap(rs.Metadata),
		Name:      rs.Name,
		Timeout:   time.Duration(rs.Timeout),
		Resources: newResources(rs.Resources),
	}
}

// ToProto converts the run settings to the protobuf representation.
func (rs *RunSettings) ToProto() *types.RunSettings {
	return &types.RunSettings{
		Command:   rs.Command,
		Image:     rs.Image,
		Queue:     rs.Queue,
		Metadata:  mkStruct(rs.Metadata),
		Name:      rs.Name,
		Timeout:   rs.Timeout.Nanoseconds(),
		Resources: rs.Resources.toProto(),
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
		rs.Resources = t.DefaultResources
	}

	return nil
}

// RepoConfig is the global configuration for the repository. It allows setting
// of global attributes as well as overrides and defaults for certain
// task-related items. It is typically named `tinyci.yml`.
type RepoConfig struct {
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
func NewRepoConfig(buf []byte) (*RepoConfig, *errors.Error) {
	r := &RepoConfig{}

	if err := yaml.UnmarshalStrict(buf, r); err != nil {
		return nil, errors.New(err)
	}

	r.handleOverrides()

	return r, r.Validate()
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
