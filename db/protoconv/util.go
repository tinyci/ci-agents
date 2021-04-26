package protoconv

import (
	"time"

	"github.com/volatiletech/null/v8"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// makeStatus returns nil if set is false and the bool is false, indicating
// that it is not set.
func makeStatus(res, set bool) null.Bool {
	if set {
		return null.BoolFrom(res)
	}

	return null.BoolFromPtr(nil)
}

func timeFromPB(ts *timestamppb.Timestamp) *time.Time {
	var t *time.Time

	if ts != nil && ts.IsValid() {
		tm := ts.AsTime()
		t = &tm
	}

	return t
}

func timeToPB(t null.Time) *timestamppb.Timestamp {
	var ret *timestamppb.Timestamp
	if t.Valid {
		ret = timestamppb.New(t.Time)
	}

	return ret
}
