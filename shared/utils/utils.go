package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func LeftPad(str string, length int, padWith string) string {
	return strings.Repeat(padWith, length-len(str)) + str
}

func formatTimePart(time int64) string {
	const length int = 2
	const base int = 10
	const padWith string = "0"
	return LeftPad(strconv.FormatInt(time, base), length, padWith)
}

// returns a time string in this format:
// days:HH:MM:SS,MsMS
func FormatTime(duration time.Duration) string {
	milllisencods := duration.Milliseconds() % 100
	seconds := int64(duration.Seconds()) % 60
	minutes := int64(duration.Minutes()) % 60
	hours := int64(duration.Hours()) % 60

	msStr := formatTimePart(milllisencods)
	secStr := formatTimePart(seconds)
	minStr := formatTimePart(minutes)
	hourStr := formatTimePart(hours)

	return fmt.Sprintf("%v:%v:%v,%v", hourStr, minStr, secStr, msStr)
}
