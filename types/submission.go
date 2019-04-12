package types

import (
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/utils"
)

// Submission is the encapsulation of a submission to the queuesvc.
type Submission struct {
	Parent      string `json:"parent"`
	Fork        string `json:"fork"`
	HeadSHA     string `json:"head_sha"`
	BaseSHA     string `json:"base_sha"`
	PullRequest int64  `json:"pull_request"`
	SubmittedBy string `json:"submitted_by"`
	All         bool   `json:"all"`

	Manual bool `json:"-"`
}

// Validate validates the submission, and returns an error if it encounters any.
func (sub *Submission) Validate() *errors.Error {
	if !sub.Manual && !utils.IsOwnerRepo(sub.Parent) {
		return errors.New("parent is invalid")
	}

	if sub.All && !sub.Manual {
		return errors.New("hook-triggered submissions may not force all")
	}

	if !utils.IsOwnerRepo(sub.Fork) {
		return errors.New("fork is invalid")
	}

	if sub.HeadSHA == "" {
		return errors.New("head sha is empty")
	}

	if !utils.IsSHA(sub.HeadSHA) {
		var err *errors.Error
		sub.HeadSHA, err = utils.QualifyBranch(sub.HeadSHA)
		if err != nil {
			return err
		}
	}

	if sub.BaseSHA == "" && !sub.Manual {
		return errors.New("base sha is empty")
	} else if !sub.Manual && !utils.IsSHA(sub.BaseSHA) {
		var err *errors.Error
		sub.BaseSHA, err = utils.QualifyBranch(sub.BaseSHA)
		if err != nil {
			return err
		}
	}

	return nil
}
