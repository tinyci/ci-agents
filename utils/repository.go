package utils

import (
	"path"
	"strings"

	"github.com/tinyci/ci-agents/errors"
)

var (
	//ErrInvalidCharacters occurs when one of "%><$." is found in either Owner or repo
	ErrInvalidCharacters = errors.New("repository name contains invalid characters")
)

// OwnerRepo returns the owner and repository parts of a github full repository
// name. It returns an *errors.Error if there are issues.
func OwnerRepo(repoName string) (string, string, *errors.Error) {
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return "", "", errors.New("parsing repository name: invalid number of parts")
	}

	for _, part := range parts {
		if len(part) == 0 {
			return "", "", errors.New("repository name part is empty")
		}

		if strings.HasPrefix(part, ".") || strings.HasSuffix(part, ".") {
			return "", "", ErrInvalidCharacters
		}
		if strings.Contains(part, ">") || strings.Contains(part, "<") {
			return "", "", ErrInvalidCharacters
		}
		if strings.Contains(part, "&") || strings.Contains(part, "%") {
			return "", "", ErrInvalidCharacters
		}
	}

	return parts[0], parts[1], nil
}

// IsOwnerRepo is a predicate to determine whether or not a github repo name is
// valid.
func IsOwnerRepo(repoName string) bool {
	_, _, err := OwnerRepo(repoName)
	return err == nil
}

// IsSHA detects if the string looks like a real SHA1 digest or not.
func IsSHA(sha string) bool {
	sha = strings.TrimSpace(sha)

	if len(sha) != 40 {
		return false
	}

	for _, c := range strings.ToLower(sha) {
		if c < '0' || c > 'f' {
			return false
		}
	}

	return true
}

// QualifyBranch corrects branch data to reflect our internal ref-type/ref-name
// strategy for tracking branches. It returns an *errors.Error when it can't figure it out.
func QualifyBranch(branch string) (string, *errors.Error) {
	if IsSHA(branch) {
		return branch, errors.New("is not a branch; is a sha")
	}

	if strings.HasPrefix(branch, "/") {
		return "", errors.New("invalid branch name")
	}

	branch = strings.TrimSpace(branch)
	branch = strings.TrimPrefix(branch, "refs/heads/")

	paths := strings.Split(branch, "/")
	if len(paths) == 0 {
		return "", errors.New("invalid branch name")
	}

	cleanedPaths := []string{}

	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p != "" && p != "." && p != ".." { // guard against relative dir shit
			cleanedPaths = append(cleanedPaths, p)
		}
	}

	if len(cleanedPaths) == 0 {
		return "", errors.New("paths were invalid")
	}

	if cleanedPaths[0] == "heads" {
		if len(cleanedPaths) > 1 {
			return path.Join(cleanedPaths...), nil
		}

		return "", errors.New("paths were invalid")
	}

	return path.Join("heads", path.Join(cleanedPaths...)), nil
}
