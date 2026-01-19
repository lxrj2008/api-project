package utils

import "math"

const (
	defaultPage = 1
	defaultSize = 20
	maxPageSize = 100
)

// NormalizePage clamps invalid pagination values.
func NormalizePage(page int) int {
	if page < 1 {
		return defaultPage
	}
	return page
}

// NormalizeSize limits the requested page size.
func NormalizeSize(size int) int {
	if size <= 0 {
		return defaultSize
	}
	return int(math.Min(float64(size), maxPageSize))
}

// Offset calculates OFFSET for SQL queries.
func Offset(page, size int) int {
	return (NormalizePage(page) - 1) * NormalizeSize(size)
}
