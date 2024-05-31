package webapi

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var formDatetimeRegex = regexp.MustCompile(
	`(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})`,
)

func ParseFormDatetime(v string, location *time.Location) (time.Time, error) {
	match := formDatetimeRegex.FindStringSubmatch(v)
	if len(match) < 1  {
		return time.Time{}, fmt.Errorf("failed to match")
	}

	year, err := strconv.Atoi(match[1])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(match[2])
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(match[3])
	if err != nil {
		return time.Time{}, err
	}

	hour, err := strconv.Atoi(match[4])
	if err != nil {
		return time.Time{}, err
	}

	minute, err := strconv.Atoi(match[5])
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		year, time.Month(month), day, hour, minute, 0, 0, location,
	).UTC(), nil
}
