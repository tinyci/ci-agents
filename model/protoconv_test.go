package model

import (
	"math/rand"
	"time"

	check "github.com/erikh/check"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func (ms *modelSuite) TestTimeProtoConversions(c *check.C) {
	epoch := time.Unix(0, 0)
	hello := time.Unix(1, 2)
	ts := MakeTimestamp(&hello)
	c.Assert(ts, check.FitsTypeOf, &timestamp.Timestamp{})
	c.Assert(ts.Seconds, check.Equals, int64(1))
	c.Assert(ts.Nanos, check.Equals, int32(2))

	// timestamps with zero value and the nullable argument as true should yield nil
	c.Assert(MakeTime(&timestamp.Timestamp{}, true), check.IsNil)
	c.Assert(MakeTime(&timestamp.Timestamp{}, false), check.DeepEquals, &epoch)
	c.Assert(MakeTime(ts, true), check.DeepEquals, &hello)

	rand.Seed(time.Now().UnixNano())
	secs := rand.Int63()
	nanos := rand.Int63n(int64(time.Second))
	r := time.Unix(int64(secs), int64(nanos))
	ts = MakeTimestamp(&r)
	c.Assert(ts, check.FitsTypeOf, &timestamp.Timestamp{})
	c.Assert(ts.Seconds, check.Equals, int64(secs), check.Commentf("%v", ts.Seconds-secs))
	c.Assert(ts.Nanos, check.Equals, int32(nanos))
}

func (ms *modelSuite) TestStatusProtoConversions(c *check.C) {
	c.Assert(MakeStatus(true, false), check.IsNil)
	c.Assert(MakeStatus(false, false), check.IsNil)
	c.Assert(*(MakeStatus(true, true)), check.Equals, true)
	c.Assert(*(MakeStatus(false, true)), check.Equals, false)
}
