package model

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
)

// MakeTimestamp takes a time.Time and makes a protobuf timestamp out of it.
func MakeTimestamp(time *time.Time) *timestamp.Timestamp {
	stamp := &timestamp.Timestamp{}
	if time != nil {
		stamp.Seconds = time.Unix()
		stamp.Nanos = int32(time.UnixNano())
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

// MakeUsers converts a proto userlist to a model one.
func MakeUsers(users []*types.User) ([]*User, *errors.Error) {
	ret := []*User{}
	for _, user := range users {
		u, err := NewUserFromProto(user)
		if err != nil {
			return nil, err
		}
		ret = append(ret, u)
	}

	return ret, nil
}

// MakeUserList returns the inverse of MakeUsers.
func MakeUserList(users []*User) []*types.User {
	ret := []*types.User{}
	for _, user := range users {
		ret = append(ret, user.ToProto())
	}

	return ret
}

// MakeStatus returns nil if set is false and the bool is false, indicating
// that it is not set.
func MakeStatus(res, set bool) *bool {
	if set {
		return &res
	}

	return nil
}
