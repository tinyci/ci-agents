package model

import (
	"fmt"

	"github.com/erikh/check"
	"github.com/golang/mock/gomock"
	"github.com/tinyci/ci-agents/mocks/github"
	"github.com/tinyci/ci-agents/utils"
)

func (ms *modelSuite) TestCancellationByRef(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	qis, err := ms.FillQueue(1000)
	c.Assert(err, check.IsNil)

	for _, qi := range qis {
		owner, repo, err := qi.Run.Task.Submission.BaseRef.Repository.OwnerRepo()
		c.Assert(err, check.IsNil)

		client.EXPECT().ErrorStatus(gomock.Any(), owner, repo, qi.Run.Name, qi.Run.Task.Submission.HeadRef.SHA, fmt.Sprintf("__test__/log/%d", qi.Run.ID), utils.ErrRunCanceled).Return(nil)
		c.Assert(ms.model.CancelRefByName(qi.Run.Task.Submission.HeadRef.Repository.ID, qi.Run.Task.Submission.HeadRef.RefName, "__test__", client), check.IsNil)

		runs, err := ms.model.GetRunsForTask(qi.Run.Task.ID, 0, 100)
		c.Assert(err, check.IsNil)

		for _, run := range runs {
			c.Assert(run.Task.Canceled, check.Equals, true)
		}
	}
}

func (ms *modelSuite) TestCancellationByTask(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	qis, err := ms.FillQueue(1000)
	c.Assert(err, check.IsNil)

	for _, qi := range qis {
		owner, repo, err := qi.Run.Task.Submission.BaseRef.Repository.OwnerRepo()
		c.Assert(err, check.IsNil)

		client.EXPECT().ErrorStatus(gomock.Any(), owner, repo, qi.Run.Name, qi.Run.Task.Submission.HeadRef.SHA, fmt.Sprintf("__test__/log/%d", qi.Run.ID), utils.ErrRunCanceled).Return(nil)
		c.Assert(ms.model.CancelTask(qi.Run.Task, "__test__", client), check.IsNil)

		runs, err := ms.model.GetRunsForTask(qi.Run.Task.ID, 0, 100)
		c.Assert(err, check.IsNil)

		for _, run := range runs {
			c.Assert(run.Task.Canceled, check.Equals, true)
		}
	}
}

func (ms *modelSuite) TestCancellationByRun(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	qis, err := ms.FillQueue(1000)
	c.Assert(err, check.IsNil)

	for _, qi := range qis {
		owner, repo, err := qi.Run.Task.Submission.BaseRef.Repository.OwnerRepo()
		c.Assert(err, check.IsNil)

		client.EXPECT().ErrorStatus(gomock.Any(), owner, repo, qi.Run.Name, qi.Run.Task.Submission.HeadRef.SHA, fmt.Sprintf("__test__/log/%d", qi.Run.ID), utils.ErrRunCanceled).Return(nil)
		c.Assert(ms.model.CancelRun(qi.Run.ID, "__test__", client), check.IsNil)

		runs, err := ms.model.GetRunsForTask(qi.Run.Task.ID, 0, 100)
		c.Assert(err, check.IsNil)

		for _, run := range runs {
			c.Assert(run.Task.Canceled, check.Equals, true)
		}
	}
}
