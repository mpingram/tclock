package main

import (
	"fmt"
	"time"
)

func durFormat(dur time.Duration) string {
	var hours, minutes int64
	hours = int64(dur.Hours())
	minutes = int64(dur.Minutes()) % 60
	return fmt.Sprintf("%vh : %vm", hours, minutes)
}
