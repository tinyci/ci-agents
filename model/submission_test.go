package model

import (
	"strings"

	check "github.com/erikh/check"
	"github.com/golang/mock/gomock"
	"github.com/tinyci/ci-agents/mocks/github"
	"github.com/tinyci/ci-agents/types"
)

var Fixtures = map[bool][]*types.Submission{
	true: {
		{
			Parent:  "foo/bar",
			Fork:    "bar/foo",
			BaseSHA: "f0d22d94df0f45a1fff37e9cd8772e7a6c2439b1",
			HeadSHA: "00c60ef6bd2cc54680205c7f5ad6639540e15cee",
		},
		{
			Parent:  "foo/bar",
			Fork:    "bar/foo",
			BaseSHA: "22cb110a32c3573250f0e6e544ad12986b31579d",
			HeadSHA: "6692fe9c58867dab715f065786b02f7146a597ce",
		},
		{
			Parent:  "foo.bar/bar.foo",
			Fork:    "bar/foo",
			BaseSHA: "4f1b10fbd4e5c2d8b331a93f3a28594e507d01bc",
			HeadSHA: "4379f1091ecdb5a4a630d0d7ea4b3137758285d1",
		},
		{
			Parent:  "foo.bar/bar.foo",
			Fork:    "bar/foo",
			BaseSHA: "97bcd1cb2b075d1bf5d1883a83cfdd6d5efbae74",
			HeadSHA: "64113585931932a97e60fa0f7c319b5a9172adf8",
			All:     true,
			Manual:  true,
		},
	},
	false: {
		{
			Parent:  "/",
			Fork:    "bar/foo",
			BaseSHA: "1fa1b2fbb2038847474aa7677957b960e8e7764e",
			HeadSHA: "69a553b699c7dc72087a7bda182fd1d9c224c5fc",
			All:     true,
		},
		{
			Parent:  "/",
			Fork:    "bar/foo",
			BaseSHA: "f00b7018674675f9260156ecbeb101f7c330877d",
			HeadSHA: "0024dc6d72f440b4cd156e242d9ebe37d4ec9ceb",
		},
		{
			Parent:  "../",
			Fork:    "bar/foo",
			BaseSHA: "51b950a56208d4fdfaae6a835de981307a9a4581",
			HeadSHA: "94f15bab50ad84fa547932811b04de2618c7ce36",
		},
		{
			Parent:  "/..",
			Fork:    "bar/foo",
			BaseSHA: "d78e460035d6f546044f43e537387432981cadb6",
			HeadSHA: "d1926d6b9e5dbe32e6a9353b5274e3936d6f36e8",
		},
		{
			Parent:  "../..",
			Fork:    "bar/foo",
			BaseSHA: "ab07b817b8c094ea8b58704ec327a8331851becd",
			HeadSHA: "436977d8685346c19e946e9ee755ed1afe789fb7",
		},
		{
			Parent:  "./.",
			Fork:    "bar/foo",
			BaseSHA: "d1c6cb0ff700a6f773daae258740a399c7937580",
			HeadSHA: "b9f0f5fed47b773c555f921e3251735dde249746",
		},
		{
			Parent:  "/.",
			Fork:    "bar/foo",
			BaseSHA: "04924728421ac094bef214a51efa19eec9586110",
			HeadSHA: "6be1c07c370c7ec719d485478208f2f4f8f6d2f2",
		},
		{
			Parent:  "./",
			Fork:    "bar/foo",
			BaseSHA: "5c5c71f4631b2334a703ca0ce8f936d8fcfe2ede",
			HeadSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
		{
			Parent:  "bar/foo/",
			Fork:    "bar/foo",
			BaseSHA: "5c5c71f4631b2334a703ca0ce8f936d8fcfe2ede",
			HeadSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
		{
			Parent:  "",
			Fork:    "bar/foo",
			BaseSHA: "5c5c71f4631b2334a703ca0ce8f936d8fcfe2ede",
			HeadSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
		{
			Parent:  "bar/foo",
			Fork:    "/",
			BaseSHA: "5c5c71f4631b2334a703ca0ce8f936d8fcfe2ede",
			HeadSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
		{
			Parent:  "bar/foo",
			Fork:    "bar/foo/",
			BaseSHA: "5c5c71f4631b2334a703ca0ce8f936d8fcfe2ede",
			HeadSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
		{
			Parent:  "bar/foo",
			Fork:    "",
			BaseSHA: "5c5c71f4631b2334a703ca0ce8f936d8fcfe2ede",
			HeadSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
		{
			Parent:  "",
			Fork:    "bar/foo",
			BaseSHA: "5c5c71f4631b2334a703ca0ce8f936d8fcfe2ede",
			HeadSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
		{
			Parent:  "",
			Fork:    "bar/foo",
			HeadSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
		{
			Parent:  "",
			Fork:    "bar/foo",
			BaseSHA: "98a974078c07b903beca23a9798f4bf03b6854d2",
		},
	},
}

func (ms *modelSuite) TestSubmissionEntries(c *check.C) {
	for success, subs := range Fixtures {
		for _, sub := range subs {
			if success {
				s, err := ms.CreateSubmission(sub)
				c.Assert(err, check.IsNil)

				c.Assert(s.ID, check.Not(check.Equals), int64(0), check.Commentf("%v", sub))
				ps := s.ToProto()

				c.Assert(ps.Id, check.Equals, s.ID)

				if ps.User != nil {
					c.Assert(ps.User.Id, check.Equals, s.User.ID)
				}

				c.Assert(ps.BaseRef.Id, check.Equals, s.BaseRef.ID)

				if ps.HeadRef != nil {
					c.Assert(ps.HeadRef.Id, check.Equals, s.HeadRef.ID)
				}

				s2, err := NewSubmissionFromProto(ps)
				c.Assert(err, check.IsNil)

				c.Assert(s.ID, check.Equals, s2.ID)

				if s2.User != nil {
					c.Assert(s.User.ID, check.Equals, s2.User.ID)
				}

				c.Assert(s.BaseRef.ID, check.Equals, s2.BaseRef.ID)

				if s2.HeadRef != nil {
					c.Assert(s.HeadRef.ID, check.Equals, s2.HeadRef.ID)
				}
			} else {
				_, err := ms.model.NewSubmissionFromMessage(sub)
				c.Assert(err, check.NotNil)
			}
		}
	}

	for success, subs := range Fixtures {
		if success {
			s, err := ms.model.SubmissionList(0, int64(len(subs)), "", "")
			c.Assert(err, check.IsNil)
			c.Assert(len(s), check.Equals, len(subs))

			card := map[string][]*types.Submission{}
			for _, sub := range subs {
				card[sub.Parent] = append(card[sub.Parent], sub)
			}

			for parent, subs := range card {
				s, err := ms.model.SubmissionList(0, 100, parent, "")
				c.Assert(err, check.IsNil)
				c.Assert(len(s), check.Not(check.Equals), 0)
				c.Assert(len(s), check.Equals, len(subs))

				s, err = ms.model.SubmissionListForRepository(parent, "", 0, 100)
				c.Assert(err, check.IsNil)
				c.Assert(len(s), check.Not(check.Equals), 0)
				c.Assert(len(s), check.Equals, len(subs))

				count, err := ms.model.SubmissionCount(parent, "")
				c.Assert(err, check.IsNil, check.Commentf("%v", parent))
				c.Assert(count, check.Equals, int64(len(subs)))
			}

			card = map[string][]*types.Submission{}
			for _, sub := range subs {
				key := strings.Join([]string{sub.Parent, sub.BaseSHA}, "_")
				card[key] = append(card[key], sub)
			}

			for _, subs := range card {
				s, err := ms.model.SubmissionList(0, 100, subs[0].Parent, subs[0].BaseSHA)
				c.Assert(err, check.IsNil)
				c.Assert(len(s), check.Not(check.Equals), 0)
				c.Assert(len(s), check.Equals, len(subs), check.Commentf("parent: %v, sha: %v", subs[0].Parent, subs[0].BaseSHA))

				s, err = ms.model.SubmissionListForRepository(subs[0].Parent, subs[0].BaseSHA, 0, 100)
				c.Assert(err, check.IsNil)
				c.Assert(len(s), check.Not(check.Equals), 0)
				c.Assert(len(s), check.Equals, len(subs))

				count, err := ms.model.SubmissionCount(subs[0].Parent, subs[0].BaseSHA)
				c.Assert(err, check.IsNil)
				c.Assert(count, check.Equals, int64(1))
			}

			count, err := ms.model.SubmissionCount("", "")
			c.Assert(err, check.IsNil)
			c.Assert(count, check.Equals, int64(len(subs)))
		}
	}
}

func (ms *modelSuite) TestSubmissionTasks(c *check.C) {
	mock := gomock.NewController(c)
	for success, subs := range Fixtures {
		if success {
			subTaskMap := map[int64][]*Task{}
			subIDMap := map[int64]*Submission{}

			for _, sub := range subs {
				s, err := ms.CreateSubmission(sub)
				c.Assert(err, check.IsNil)

				subIDMap[s.ID] = s

				for i := 0; i < 200; i++ {
					t, err := ms.CreateTaskForSubmission(s)
					c.Assert(err, check.IsNil)

					subTaskMap[s.ID] = append(subTaskMap[s.ID], t)
				}
			}

			for subID, tasks := range subTaskMap {
				s, err := ms.model.GetSubmissionByID(subID)
				c.Assert(err, check.IsNil)
				c.Assert(s.ID, check.Equals, subID)
				c.Assert(s.TasksCount, check.Equals, int64(len(tasks)))

				// reverse the list since the ordering will be returned in rev chrono order
				t2 := []*Task{}

				for i := len(tasks) - 1; i >= 0; i-- {
					t2 = append(t2, tasks[i])
				}

				for i := int64(0); i < 2; i++ {
					ts, err := ms.model.TasksForSubmission(s, i, 100)
					c.Assert(err, check.IsNil)
					c.Assert(len(ts), check.Equals, 100)

					for x, task := range ts {
						c.Assert(t2[int64(x)+(i*100)].ID, check.Equals, task.ID, check.Commentf("%v", int64(x)+(i*100)))
					}
				}

				gh := github.NewMockClient(mock)
				var runCount int64
				runs := []*Run{}

				for _, task := range tasks {
					owner, repo, err := task.Submission.BaseRef.Repository.OwnerRepo()
					c.Assert(err, check.IsNil)
					count, err := ms.model.CountRunsForTask(task.ID)
					c.Assert(err, check.IsNil)

					var i int64
					for {
						r, err := ms.model.GetRunsForTask(task.ID, i, 100)
						c.Assert(err, check.IsNil)
						runs = append(runs, r...)
						if len(r) == 0 {
							break
						}
						i++
					}

					runCount += count

					for i := int64(0); i < count; i++ {
						gh.EXPECT().ErrorStatus(gomock.Any(), owner, repo, "default", task.Submission.HeadRef.SHA, gomock.Any(), gomock.Any())
					}
				}

				c.Assert(s.RunsCount, check.Equals, runCount)

				newRuns := []*Run{}
				var i int64
				for {
					r, err := ms.model.RunsForSubmission(subIDMap[subID], i, 100)
					c.Assert(err, check.IsNil)

					newRuns = append(newRuns, r...)
					if len(r) == 0 {
						break
					}
					i++
				}

				c.Assert(runs, check.DeepEquals, newRuns)

				// checking that we can cancel submissions with canceled tasks in them
				c.Assert(ms.model.CancelTask(tasks[0], "", gh), check.IsNil)
				c.Assert(ms.model.CancelSubmissionByID(subID, "", gh), check.IsNil)
				s, err = ms.model.GetSubmissionByID(subID)
				c.Assert(err, check.IsNil)
				c.Assert(s.Canceled, check.Equals, true)

				for i := int64(0); i < 2; i++ {
					ts, err := ms.model.TasksForSubmission(s, i, 100)
					c.Assert(err, check.IsNil)
					c.Assert(len(ts), check.Equals, 100)

					for _, task := range ts {
						c.Assert(task.Canceled, check.Equals, true)
					}
				}
			}
		}
	}
}
