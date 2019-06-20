package model

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// MakeTimestamp takes a time.Time and makes a protobuf timestamp out of it.
func MakeTimestamp(t *time.Time) *timestamp.Timestamp {
	stamp := &timestamp.Timestamp{}

	if t != nil {
		stamp.Seconds = t.Unix()
		nanos := t.UnixNano() - t.Unix()*int64(time.Second)

		if nanos < 0 || nanos > int64(time.Second) {
			// cannot make a nanosecond fraction longer than a second or less than zero
			stamp.Nanos = 0
		} else {
			stamp.Nanos = int32(nanos)
		}
	}

	return stamp
}

// MakeTime takes a protobuf timestamp and makes a time.Time out of it. If you
// pass true to the second argument, it will return nil if the argument is the
// zero value. Otherwise, it returns a time.Unix(0, 0).
func MakeTime(ts *timestamp.Timestamp, nullable bool) *time.Time {
	if ts.Seconds != 0 || ts.Nanos != 0 {
		u := time.Unix(ts.Seconds, int64(ts.Nanos))
		return &u
	}

	if nullable {
		return nil
	}

	t := time.Unix(0, 0)
	return &t
}

// MakeStatus returns nil if set is false and the bool is false, indicating
// that it is not set.
func MakeStatus(res, set bool) *bool {
	if set {
		return &res
	}

	return nil
}
