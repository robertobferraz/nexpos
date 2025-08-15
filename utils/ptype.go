package utils

import (
	"time"
)

func PString(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func PFloat64(f float64) *float64 {
	return &f
}

func PInt(i int) *int {
	return &i
}

func PTime(t time.Time) *time.Time {
	return &t
}

func PDuration(d time.Duration) *time.Duration {
	return &d
}

func PBool(b bool) *bool {
	return &b
}

func PByte(b []byte) *[]byte {
	return &b
}

func PInt64(i int64) *int64 {
	return &i
}
