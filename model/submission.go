package model

import (
	"time"

	gtypes "github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/types"
)

// Submission is the concrete type for a test submission, unlike
// types/submission.go which is the ephemeral type for messages and so on.
// Each submission has it own ID and list of tasks so that we can further give
// the user the opportunity to aggregate their tasks into a single unit.
type Submission struct {
	ID int64 `gorm:"priamry_key" json:"id"`

	User   *User `gorm:"association_autoupdate:false,nullable:true" json:"user"`
	UserID int64 `json:"-"`

	HeadRef   *Ref  `gorm:"association_autoupdate:false,column:head_ref_id,nullable:true" json:"head_ref"`
	HeadRefID int64 `json:"-"`

	BaseRef   *Ref  `gorm:"association_autoupdate:false,column:base_ref_id" json:"base_ref"`
	BaseRefID int64 `json:"-"`

	CreatedAt time.Time
}

// ToProto converts the submissions to the protobuf version
func (s *Submission) ToProto() *gtypes.Submission {
	var (
		pu *gtypes.User
		hr *gtypes.Ref
	)

	if s.User != nil {
		pu = s.User.ToProto()
	}

	if s.HeadRef != nil {
		hr = s.HeadRef.ToProto()
	}

	return &gtypes.Submission{
		Id:        s.ID,
		BaseRef:   s.BaseRef.ToProto(),
		HeadRef:   hr,
		User:      pu,
		CreatedAt: MakeTimestamp(&s.CreatedAt),
	}
}

// NewSubmissionFromProto converts the proto representation to the task type.
func NewSubmissionFromProto(gt *gtypes.Submission) (*Submission, *errors.Error) {
	var (
		u       *User
		headref *Ref
		err     *errors.Error
	)

	if gt.User != nil {
		u, err = NewUserFromProto(gt.User)
		if err != nil {
			return nil, err.Wrap("converting for use in submission")
		}
	}

	if gt.HeadRef != nil {
		headref, err = NewRefFromProto(gt.HeadRef)
		if err != nil {
			return nil, err.Wrap("converting for use in submission")
		}
	}

	baseref, err := NewRefFromProto(gt.BaseRef)
	if err != nil {
		return nil, err.Wrap("converting for use in submission")
	}

	return &Submission{
		ID:      gt.Id,
		User:    u,
		BaseRef: baseref,
		HeadRef: headref,
	}, nil
}

// NewSubmissionFromMessage creates a blank record populated by the appropriate data
// reaped from the message type in types/submission.go.
func (m *Model) NewSubmissionFromMessage(sub *types.Submission) (*Submission, *errors.Error) {
	if err := sub.Validate(); err != nil {
		return nil, err.Wrap("did not pass validation")
	}

	var (
		u    *User
		base *Ref
		err  *errors.Error
	)

	if sub.SubmittedBy != "" {
		u, err = m.FindUserByName(sub.SubmittedBy)
		if err != nil {
			return nil, err.Wrap("for use in submission")
		}
	}

	if sub.BaseSHA != "" {
		base, err = m.GetRefByNameAndSHA(sub.Parent, sub.BaseSHA)
		if err != nil {
			return nil, err.Wrap("base")
		}
	}

	head, err := m.GetRefByNameAndSHA(sub.Fork, sub.HeadSHA)
	if err != nil {
		return nil, err.Wrap("head")
	}

	return &Submission{
		User:    u,
		HeadRef: head,
		BaseRef: base,
	}, nil
}
