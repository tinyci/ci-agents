package utils

import (
	"strconv"

	"errors"
)

// MaxPerPage is the maximum size of a single page that tinyCI will generate when rendering responses.
const MaxPerPage int64 = 100

const defaultPerPage int64 = 10

func intp(p int64) *int64 {
	return &p
}

// ScopePaginationInt applies constraints to pagination; it sets a maximum and a
// default if not supplied to perPage. Integer-only version.
func ScopePaginationInt(page, perPage *int64) (int, int, error) {
	// FIXME make this int32 everywhere
	if page == nil {
		page = intp(0)
	}

	if perPage == nil {
		perPage = intp(defaultPerPage)
	}

	if *perPage == 0 {
		*perPage = defaultPerPage
	}

	if *page < 0 {
		return 0, 0, errors.New("invalid page")
	}

	if *perPage > MaxPerPage {
		*perPage = MaxPerPage
	}

	if *perPage < 0 {
		return 0, 0, errors.New("invalid per page")
	}

	return int(*page), int(*perPage), nil
}

// ScopePagination applies constraints to pagination; it sets a maximum and a
// default if not supplied to perPage. String version.
func ScopePagination(pg, ppg string) (int, int, error) {
	var (
		page, perPage int64
	)

	if pg == "" {
		page = 0
	} else {
		p, err := strconv.ParseInt(pg, 10, 64)
		if err != nil {
			return 0, 0, err
		}

		page = p
	}

	if ppg == "" || ppg == "0" {
		perPage = defaultPerPage
	} else {
		p, err := strconv.ParseInt(ppg, 10, 64)
		if err != nil {
			return 0, 0, err
		}

		perPage = p
	}

	return ScopePaginationInt(&page, &perPage)
}
