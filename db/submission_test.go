package db

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"

	gtypes "github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/db/protoconv"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
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

func TestSubmissionEntries(t *testing.T) {
	m := testInit(t)
	c := protoconv.New(m.db)

	for success, subs := range Fixtures {
		for _, sub := range subs {
			if success {
				s, err := m.CreateTestSubmission(ctx, sub)
				assert.NilError(t, err)

				assert.Assert(t, s.ID != 0)

				sp, err := c.ToProto(ctx, s)
				assert.NilError(t, err)

				ps := sp.(*gtypes.Submission)

				assert.Assert(t, cmp.Equal(ps.Id, s.ID))

				if ps.User != nil {
					assert.Assert(t, cmp.Equal(ps.User.Id, s.UserID))
				}

				assert.Assert(t, cmp.Equal(ps.BaseRef.Id, s.BaseRefID))

				if ps.HeadRef != nil {
					assert.Assert(t, cmp.Equal(ps.HeadRef.Id, s.HeadRefID.Int64))
				}

				s2c, err := c.FromProto(ctx, ps)
				assert.NilError(t, err)

				s2 := s2c.(*models.Submission)

				assert.Assert(t, cmp.Equal(s.ID, s2.ID))
				assert.Assert(t, cmp.Equal(s.UserID, s2.UserID))
				assert.Assert(t, cmp.Equal(s.BaseRefID, s2.BaseRefID))
				assert.Assert(t, cmp.Equal(s.HeadRefID, s2.HeadRefID))
			} else {
				_, err := m.CreateTestSubmission(ctx, sub)
				assert.Assert(t, err != nil)
			}
		}
	}

	for success, subs := range Fixtures {
		if success {
			s, err := m.SubmissionList(ctx, 0, int64(len(subs)), "", "")
			assert.NilError(t, err)
			assert.Assert(t, cmp.Equal(len(s), len(subs)))

			card := map[string][]*types.Submission{}
			for _, sub := range subs {
				card[sub.Parent] = append(card[sub.Parent], sub)
			}

			for parent, subs := range card {
				s, err := m.SubmissionList(ctx, 0, 0, parent, "")
				assert.NilError(t, err)
				assert.Assert(t, len(s) != 0)
				assert.Assert(t, cmp.Equal(len(s), len(subs)))

				count, err := m.SubmissionCount(ctx, parent, "")
				assert.NilError(t, err)
				assert.Assert(t, cmp.Equal(count, int64(len(subs))), parent)
			}

			card = map[string][]*types.Submission{}
			for _, sub := range subs {
				key := strings.Join([]string{sub.Parent, sub.BaseSHA}, "_")
				card[key] = append(card[key], sub)
			}

			for key, subs := range card {
				s, err := m.SubmissionList(ctx, 0, utils.MaxPerPage, subs[0].Parent, subs[0].BaseSHA)
				assert.NilError(t, err, key)
				assert.Assert(t, len(s) != 0, key)
				assert.Assert(t, cmp.Equal(len(s), len(subs)), key)

				count, err := m.SubmissionCount(ctx, subs[0].Parent, subs[0].BaseSHA)
				assert.NilError(t, err, key)
				assert.Assert(t, cmp.Equal(count, int64(1)), key)
			}

			count, err := m.SubmissionCount(ctx, "", "")
			assert.NilError(t, err)
			assert.Assert(t, cmp.Equal(count, int64(len(subs))))
		}
	}
}

func TestSubmissionTasks(t *testing.T) {
	m := testInit(t)

	for success, subs := range Fixtures {
		if success {
			subTaskMap := map[int64][]*models.Task{}
			subIDMap := map[int64]*models.Submission{}

			for _, sub := range subs {
				s, err := m.CreateTestSubmission(ctx, sub)
				assert.NilError(t, err)

				subIDMap[s.ID] = s

				for i := int64(0); i < utils.MaxPerPage*2; i++ { // this 2 relates to the task count below in TasksForSubmission
					task, err := m.CreateTestTaskForSubmission(ctx, s)
					assert.NilError(t, err)

					subTaskMap[s.ID] = append(subTaskMap[s.ID], task)
				}
			}

			for subID, tasks := range subTaskMap {
				s, err := m.GetSubmissionByID(ctx, subID)
				assert.NilError(t, err)
				assert.Assert(t, cmp.Equal(s.ID, subID))
				count, err := s.Tasks().Count(ctx, m.db)
				assert.NilError(t, err)
				assert.Assert(t, cmp.Equal(count, int64(len(tasks))))

				for i := int64(0); i < 2; i++ {
					ts, err := m.TasksForSubmission(ctx, s.ID, i, utils.MaxPerPage)
					assert.NilError(t, err)
					assert.Assert(t, cmp.Equal(len(ts), int(utils.MaxPerPage)), i)

					for x, task := range ts {
						assert.Assert(t, cmp.Equal(tasks[int64(x)+(i*utils.MaxPerPage)].ID, task.ID), fmt.Sprintf("%v", int64(x)+(i*utils.MaxPerPage)))
					}
				}

				var runCount int64
				runs := []*models.Run{}

				for _, task := range tasks {
					count, err := m.CountRunsForTask(ctx, task.ID)
					assert.NilError(t, err)

					var i int64
					for {
						r, err := m.GetRunsForTask(ctx, task.ID, i, utils.MaxPerPage)
						assert.NilError(t, err)
						runs = append(runs, r...)
						if len(r) == 0 {
							break
						}
						i++
					}

					runCount += count
				}

				sqlRunCount, err := m.CountRunsForSubmission(ctx, subID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Assert(t, cmp.Equal(sqlRunCount, runCount))

				newRuns := []*models.Run{}
				var i int64
				for {
					r, err := m.RunsForSubmission(ctx, subID, i, utils.MaxPerPage)
					assert.NilError(t, err)

					newRuns = append(newRuns, r...)
					if len(r) == 0 {
						break
					}
					i++
				}

				assert.Assert(t, cmp.DeepEqual(runs, newRuns))

				// checking that we can cancel submissions with canceled tasks in them
				assert.NilError(t, m.CancelTask(ctx, tasks[0].ID))
				assert.NilError(t, m.CancelSubmissionByID(ctx, subID))
				s, err = m.GetSubmissionByID(ctx, subID)
				assert.NilError(t, err)

				for i := int64(0); i < 2; i++ {
					ts, err := m.TasksForSubmission(ctx, s.ID, i, utils.MaxPerPage)
					assert.NilError(t, err)
					assert.Assert(t, cmp.Equal(len(ts), int(utils.MaxPerPage)))

					for _, task := range ts {
						assert.Assert(t, task.Canceled)
					}
				}
			}
		}
	}
}
