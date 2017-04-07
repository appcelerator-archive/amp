package cli

import (
	"strconv"
	"time"
)

func ConvertTime(c Interface, in int64) time.Time {
	intVal, err := strconv.ParseInt(strconv.FormatInt(in, 10), 10, 64)
	if err != nil {
		c.Console().Fatalf(err.Error())
	}
	out := time.Unix(intVal, 0)
	return out
}
