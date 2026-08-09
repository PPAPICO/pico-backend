package parser

import (
	"regexp"
	"strings"
	"time"
)

func ParseYYYYMMDD(s string) time.Time {
	s = strings.TrimSpace(s)
	if len(s) < 8 {
		return time.Time{}
	}
	re := regexp.MustCompile(`\d{8}`)
	match := re.FindString(s)
	if match == "" {
		return time.Time{}
	}

	t, err := time.Parse("20060102", match)
	if err != nil {
		return time.Time{}
	}
	return t
}

func ParseDateRange(rangeStr string) (time.Time, time.Time) {
	if rangeStr == "" {
		return time.Time{}, time.Time{}
	}

	re := regexp.MustCompile(`\d{4}[.-/]?\d{2}[.-/]?\d{2}`)
	matches := re.FindAllString(rangeStr, -1)

	var startDate, endDate time.Time

	if len(matches) >= 1 {
		startDate = cleanAndParseDate(matches[0])
	}
	if len(matches) >= 2 {
		endDate = cleanAndParseDate(matches[1])
	}

	return startDate, endDate
}

func cleanAndParseDate(dateStr string) time.Time {
	cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(dateStr, "")
	if len(cleaned) == 8 {
		t, err := time.Parse("20060102", cleaned)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}
