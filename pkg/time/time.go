package time

import (
	"time"
)

const (
	humanFormat = "02 Jan 06 15:04"
)

// ConvertTime converts unix timestamp to human readable time
func ConvertTime(in int64) string {
	return time.Unix(in, 0).Format(humanFormat)
}
