package conversion

import "time"

func PointerToString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func PointerToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return t.UTC()
}
