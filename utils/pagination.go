package utils

import (
	"strconv"

	"github.com/tinyci/ci-agents/errors"
)

// MaxPerPage is the maximum size of a single page that tinyCI will generate when rendering responses.
const MaxPerPage int64 = 100

const defaultPerPage int64 = 10

// ScopePaginationInt applies constraints to pagination; it sets a maximum and a
// default if not supplied to perPage. Integer-only version.
func ScopePaginationInt(page, perPage int64) (int64, int64, error) {
	if page < 0 {
		return 0, 0, errors.New("invalid page")
	}

	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	if perPage < 0 {
		return 0, 0, errors.New("invalid per page")
	}

	if perPage == 0 {
		perPage = defaultPerPage
	}

	return page, perPage, nil
}

// ScopePagination applies constraints to pagination; it sets a maximum and a
// default if not supplied to perPage. String version.
func ScopePagination(pg, ppg string) (int64, int64, error) {
	var (
		err           error
		page, perPage int64
	)

	if pg == "" {
		page = 0
	} else {
		page, err = strconv.ParseInt(pg, 10, 64)
		if err != nil {
			return 0, 0, errors.New(err)
		}
	}

	if ppg == "" || ppg == "0" {
		perPage = defaultPerPage
	} else {
		perPage, err = strconv.ParseInt(ppg, 10, 64)
		if err != nil {
			return 0, 0, errors.New(err)
		}
	}

	return ScopePaginationInt(page, perPage)
}
