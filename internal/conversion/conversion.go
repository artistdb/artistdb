package conversion

import (
	"time"
)

func String(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func Time(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return t.UTC()
}

func TimeP(t time.Time) *time.Time {
	return &t
}

func RFC3339P(t time.Time) *string {
	if t.IsZero() {
		return nil
	}

	ret := t.Format(time.RFC3339)
	return &ret
}
