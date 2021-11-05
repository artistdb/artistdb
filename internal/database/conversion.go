package database

import "time"

func toString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func toTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return t.UTC()
}
