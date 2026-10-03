package cmd

import (
	"fmt"
	"time"
)

func formatDuration(d time.Duration) string {
	totalMinutes := int(d.Minutes())

	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %02dm", hours, minutes)
	}

	return fmt.Sprintf("%dm", minutes)
}

func progressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}

	if percent > 100 {
		percent = 100
	}

	filled := int((percent / 100) * float64(width))

	result := ""

	for i := 0; i < width; i++ {
		if i < filled {
			result += "█"
		} else {
			result += "░"
		}
	}

	return result
}