package types

import (
	"fmt"
	"io/ioutil"
	. "testing"
	"time"

	check "github.com/erikh/check" // check has a Run type or method too; so we can't dot import it
)

type typesSuite struct{}

var _ = check.Suite(&typesSuite{})

func TestTypes(t *T) {
	check.TestingT(t)
}

func (ts *typesSuite) TestTaskSettingsNew(c *check.C) {
	iters := map[string]*TaskSettings{
		"basic": {
			Mountpoint:   "/tmp",
			WorkDir:      "/foobar",
			DefaultQueue: "frobnik",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "frobnik",
					Name:    "foobar",
				},
			},
		},

		"independent_queues": {
			Mountpoint:   "/tmp",
			WorkDir:      "/foobar",
			DefaultQueue: "frobnik",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "frobnik",
					Name:    "foobar",
				},
				"foobar2": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "hi",
					Name:    "foobar2",
				},
				"foobar3": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "hi2",
					Name:    "foobar3",
				},
			},
		},

		"standard_defaults": {
			Mountpoint: "/tmp",
			WorkDir:    "/tmp",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "default",
					Name:    "foobar",
				},
			},
		},
	}

	for file, task := range iters {
		buf, err := ioutil.ReadFile(fmt.Sprintf("testdata/tasks/%s.yml", file))
		c.Assert(err, check.IsNil)

		t, err := NewTaskSettings(buf, true, RepoConfig{})
		c.Assert(err, check.IsNil)

		c.Assert(t, check.DeepEquals, task, check.Commentf("%s", file))
	}
}

func (ts *typesSuite) TestRepoConfig(c *check.C) {
	iters := map[string]RepoConfig{
		"basic": {
			Queue: "default",
		},

		"complicated": {
			Queue:         "quux",
			WorkDir:       "/",
			OverrideQueue: true,
		},
	}

	for file, config := range iters {
		buf, err := ioutil.ReadFile(fmt.Sprintf("testdata/repo_configs/%s.yml", file))
		c.Assert(err, check.IsNil)

		rc, err := NewRepoConfig(buf)
		c.Assert(err, check.IsNil)

		c.Assert(rc, check.DeepEquals, config)
	}
}

func (ts *typesSuite) TestNewTaskWithRepoConfig(c *check.C) {
	iters := map[string]*TaskSettings{
		"basic": {
			Mountpoint:   "/tmp",
			WorkDir:      "/foobar",
			DefaultQueue: "frobnik",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "frobnik",
					Name:    "foobar",
				},
			},
			Config: RepoConfig{
				Queue: "default",
			},
		},

		"queue_default": {
			Mountpoint: "/tmp",
			WorkDir:    "/",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "quux",
					Name:    "foobar",
				},
			},
			Config: RepoConfig{
				Queue: "quux",
			},
		},

		"override_queue": {
			Mountpoint: "/tmp",
			WorkDir:    "/",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "quux",
					Name:    "foobar",
				},
			},
			Config: RepoConfig{
				Queue:         "quux",
				OverrideQueue: true,
			},
		},

		"timeout": {
			Mountpoint:     "/tmp",
			DefaultTimeout: 50 * time.Millisecond,
			WorkDir:        "/",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "default",
					Name:    "foobar",
					Timeout: 120 * time.Millisecond,
				},
			},
			Config: RepoConfig{
				Queue:         "default",
				GlobalTimeout: 12 * time.Millisecond,
			},
		},
		"task_default_timeout": {
			Mountpoint:     "/tmp",
			DefaultTimeout: 50 * time.Millisecond,
			WorkDir:        "/",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "default",
					Name:    "foobar",
					Timeout: 50 * time.Millisecond,
				},
			},
			Config: RepoConfig{
				Queue:         "default",
				GlobalTimeout: 12 * time.Millisecond,
			},
		},
		"repo_config_timeout": {
			Mountpoint: "/tmp",
			WorkDir:    "/",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "default",
					Name:    "foobar",
					Timeout: 12 * time.Millisecond,
				},
			},
			Config: RepoConfig{
				Queue:         "default",
				GlobalTimeout: 12 * time.Millisecond,
			},
		},
		"repo_config_timeout_override": {
			Mountpoint:     "/tmp",
			DefaultTimeout: 50 * time.Millisecond,
			WorkDir:        "/",
			Runs: map[string]*RunSettings{
				"foobar": {
					Command: []string{"foo", "bar"},
					Image:   "foobar",
					Queue:   "default",
					Name:    "foobar",
					Timeout: 12 * time.Millisecond,
				},
			},
			Config: RepoConfig{
				Queue:           "default",
				GlobalTimeout:   12 * time.Millisecond,
				OverrideTimeout: true,
			},
		},
	}

	for file, task := range iters {
		buf, err := ioutil.ReadFile(fmt.Sprintf("testdata/repo_configs/%s.yml", file))
		c.Assert(err, check.IsNil, check.Commentf("%s", file))

		rc, err := NewRepoConfig(buf)
		c.Assert(err, check.IsNil, check.Commentf("%s", file))

		buf, err = ioutil.ReadFile(fmt.Sprintf("testdata/tasks/%s.yml", file))
		c.Assert(err, check.IsNil, check.Commentf("%s", file))

		t, err := NewTaskSettings(buf, true, rc)
		c.Assert(err, check.IsNil, check.Commentf("%s", file))

		c.Assert(t, check.DeepEquals, task, check.Commentf("%s", file))
	}
}

func (ts *typesSuite) TestDefaultImage(c *check.C) {
	t := &TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*RunSettings{
			"no_image": {
				Command: []string{"cmd"},
			},
		},
		Config: RepoConfig{
			DefaultImage: "repo_config_image",
		},
	}

	t.handleOverrides()

	c.Assert(t.Runs["no_image"].Image, check.Equals, "repo_config_image")

	t = &TaskSettings{
		DefaultImage: "task_image",
		Mountpoint:   "/tmp",
		Runs: map[string]*RunSettings{
			"no_image": {
				Command: []string{"cmd"},
			},
		},
		Config: RepoConfig{
			DefaultImage: "repo_config_image",
		},
	}

	t.handleOverrides()

	c.Assert(t.Runs["no_image"].Image, check.Equals, "task_image")

	t = &TaskSettings{
		DefaultImage: "task_image",
		Mountpoint:   "/tmp",
		Runs: map[string]*RunSettings{
			"no_image": {
				Command: []string{"cmd"},
				Image:   "local_image",
			},
		},
		Config: RepoConfig{
			DefaultImage: "repo_config_image",
		},
	}

	t.handleOverrides()

	c.Assert(t.Runs["no_image"].Image, check.Equals, "local_image")
}

func (ts *typesSuite) TestRepoConfigMetadata(c *check.C) {
	t := &TaskSettings{
		Mountpoint:   "/tmp",
		WorkDir:      "/foobar",
		DefaultQueue: "frobnik",
		Runs: map[string]*RunSettings{
			"foobar": {
				Command: []string{"foo", "bar"},
				Image:   "foobar",
				Queue:   "frobnik",
				Name:    "foobar",
				Metadata: map[string]interface{}{
					"override": "yep",
				},
			},
			"baz": {
				Command: []string{"foo", "bar"},
				Image:   "foobar",
				Queue:   "frobnik",
				Name:    "foobar",
			},
		},
		Config: RepoConfig{
			Queue: "default",
			Metadata: map[string]interface{}{
				"foo":      "bar",
				"override": "nope",
			},
		},
	}

	t.handleOverrides()

	c.Assert(t.Runs["foobar"].Metadata["foo"], check.Equals, "bar")
	c.Assert(t.Runs["foobar"].Metadata["override"], check.Equals, "yep")
	c.Assert(t.Runs["baz"].Metadata["foo"], check.Equals, "bar")
	c.Assert(t.Runs["baz"].Metadata["override"], check.Equals, "nope")

	t2 := &TaskSettings{
		Mountpoint:   "/tmp",
		WorkDir:      "/foobar",
		DefaultQueue: "frobnik",
		Runs: map[string]*RunSettings{
			"foobar": {
				Command: []string{"foo", "bar"},
				Image:   "foobar",
				Queue:   "frobnik",
				Name:    "foobar",
				Metadata: map[string]interface{}{
					"override": "yep",
				},
			},
			"baz": {
				Command: []string{"foo", "bar"},
				Image:   "foobar",
				Queue:   "frobnik",
				Name:    "foobar",
			},
		},
		Config: RepoConfig{
			Queue: "default",
			Metadata: map[string]interface{}{
				"foo":      "bar",
				"override": "nope",
			},
			OverrideMetadata: true,
		},
	}

	t2.handleOverrides()

	c.Assert(t2.Runs["foobar"].Metadata["foo"], check.Equals, "bar")
	c.Assert(t2.Runs["foobar"].Metadata["override"], check.Equals, "nope")
	c.Assert(t2.Runs["baz"].Metadata["foo"], check.Equals, "bar")
	c.Assert(t2.Runs["baz"].Metadata["override"], check.Equals, "nope")
}

func (ts *typesSuite) TestTaskDefaults(c *check.C) {
	t := &TaskSettings{
		Mountpoint:     "/tmp",
		WorkDir:        "/foobar",
		DefaultImage:   "foobar",
		DefaultQueue:   "funky",
		DefaultTimeout: time.Minute,
		Runs: map[string]*RunSettings{
			"test": {
				Command: []string{"foo", "bar"},
			},
		},
	}

	c.Assert(t.Validate(true), check.IsNil)
	c.Assert(t.Runs["test"].Timeout, check.Equals, time.Minute)
	c.Assert(t.Runs["test"].Image, check.Equals, "foobar")
	c.Assert(t.Runs["test"].Queue, check.Equals, "funky")
}

func (ts *typesSuite) TestTaskDependenciesNoRuns(c *check.C) {
	t := &TaskSettings{
		Dependencies: []string{"one"},
	}

	c.Assert(t.Validate(false), check.IsNil)
	c.Assert(t.Validate(true), check.NotNil)

	content, err := ioutil.ReadFile("testdata/tasks/deps_only.yml")
	c.Assert(err, check.IsNil)

	t, err = NewTaskSettings(content, false, RepoConfig{})
	c.Assert(err, check.IsNil)
	c.Assert(t.Validate(false), check.IsNil)
}

func (ts *typesSuite) TestResourceCascade(c *check.C) {
	c.Skip("resources are currently not working")

	type predicate struct { //nolint:unused
		ts       *TaskSettings
		validate func(c *check.C, name string, t *TaskSettings)
	}

	checkCascade := func(c *check.C, name string, t *TaskSettings) {
		c.Assert(t.Validate(false), check.IsNil)

		c.Assert(t.Runs["test"].Resources, check.DeepEquals, t.DefaultResources, check.Commentf("%s", name))
		c.Assert(t.Runs["test"].Resources.CPU, check.Equals, "1", check.Commentf("%s", name))
		c.Assert(t.Runs["test"].Resources.Memory, check.Equals, "10", check.Commentf("%s", name))
		c.Assert(t.Runs["test"].Resources.Disk, check.Equals, "5", check.Commentf("%s", name))
		c.Assert(t.Runs["test"].Resources.IOPS, check.Equals, "9001", check.Commentf("%s", name))
	}

	table := map[string]predicate{
		"repoconfig-cascade": {
			ts: &TaskSettings{
				Mountpoint: "/tmp",
				WorkDir:    "/foobar",
				Runs: map[string]*RunSettings{
					"test": {
						Image:   "foo",
						Command: []string{"foo", "bar"},
					},
				},
				Config: RepoConfig{
					DefaultResources: Resources{
						CPU:    "1",
						Memory: "10",
						Disk:   "5",
						IOPS:   "9001",
					},
				},
			},
			validate: checkCascade,
		},
		"tasksetting-cascade": {
			ts: &TaskSettings{
				Mountpoint: "/tmp",
				WorkDir:    "/foobar",
				Runs: map[string]*RunSettings{
					"test": {
						Image:   "foo",
						Command: []string{"foo", "bar"},
					},
				},
				DefaultResources: Resources{
					CPU:    "1",
					Memory: "10",
					Disk:   "5",
					IOPS:   "9001",
				},
			},
			validate: checkCascade,
		},
		"tasksetting-cascade-some-runs-populated": {
			ts: &TaskSettings{
				Mountpoint: "/tmp",
				WorkDir:    "/foobar",
				Runs: map[string]*RunSettings{
					"test": {
						Image:   "foo",
						Command: []string{"foo", "bar"},
					},
					"test2": {
						Image:   "foo",
						Command: []string{"foo", "bar"},
						Resources: Resources{
							CPU:    "2",
							Memory: "10",
							Disk:   "5",
							IOPS:   "9001",
						},
					},
				},
				DefaultResources: Resources{
					CPU:    "1",
					Memory: "10",
					Disk:   "5",
					IOPS:   "9001",
				},
			},
			validate: func(c *check.C, name string, t *TaskSettings) {
				checkCascade(c, name, t)
				c.Assert(t.Runs["test2"].Resources, check.Not(check.DeepEquals), t.DefaultResources, check.Commentf("%s", name))
				c.Assert(t.Runs["test2"].Resources.CPU, check.Equals, "2", check.Commentf("%s", name))
			},
		},
		"tasksetting-run-populated": {
			ts: &TaskSettings{
				Mountpoint: "/tmp",
				WorkDir:    "/foobar",
				Runs: map[string]*RunSettings{
					"test": {
						Image:   "foo",
						Command: []string{"foo", "bar"},
						Resources: Resources{
							CPU:    "2",
							Memory: "10",
							Disk:   "5",
							IOPS:   "9001",
						},
					},
				},
				DefaultResources: Resources{
					CPU:    "1",
					Memory: "10",
					Disk:   "5",
					IOPS:   "9001",
				},
			},
			validate: func(c *check.C, name string, t *TaskSettings) {
				c.Assert(t.Validate(false), check.IsNil, check.Commentf("%s", name))
				c.Assert(t.Runs["test"].Resources, check.Not(check.DeepEquals), t.DefaultResources, check.Commentf("%s", name))
				c.Assert(t.Runs["test"].Resources.CPU, check.Equals, "2", check.Commentf("%s", name))
			},
		},
	}

	for name, pred := range table {
		pred.validate(c, name, pred.ts)
	}
}
