package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/session"
)

var monthCmd = &cobra.Command{
	Use:   "month",
	Short: "Show this month's productivity",

	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()

		year, month, _ := now.Date()

		firstDay := time.Date(
			year,
			month,
			1,
			0, 0, 0, 0,
			now.Location(),
		)

		// Day 0 of the next month = last day
		// of the current month.
		lastDay := time.Date(
			year,
			month+1,
			0,
			0, 0, 0, 0,
			now.Location(),
		)

		fmt.Println()
		fmt.Printf(
			"%s\n",
			firstDay.Format("January 2006"),
		)
		fmt.Println("============================================================")

		var total int64
		var study int64
		var work int64
		var activeDays int

		for day := 1; day <= lastDay.Day(); day++ {
			date := time.Date(
				year,
				month,
				day,
				0, 0, 0, 0,
				now.Location(),
			)

			stats, err := session.GetDailyStats(date)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			total += stats.TotalSeconds
			study += stats.StudySeconds
			work += stats.WorkSeconds

			if stats.TotalSeconds > 0 {
				activeDays++
			}

			duration := time.Duration(stats.TotalSeconds) *
				time.Second

			bar := monthBar(duration)

			fmt.Printf(
				"%s  %-8s %s\n",
				date.Format("02 Mon"),
				formatDuration(duration),
				bar,
			)
		}

		fmt.Println("------------------------------------------------------------")

		fmt.Printf(
			"Active days : %d / %d\n",
			activeDays,
			lastDay.Day(),
		)

		fmt.Printf(
			"Study       : %s\n",
			formatDuration(
				time.Duration(study)*time.Second,
			),
		)

		fmt.Printf(
			"Work        : %s\n",
			formatDuration(
				time.Duration(work)*time.Second,
			),
		)

		fmt.Printf(
			"Total       : %s\n",
			formatDuration(
				time.Duration(total)*time.Second,
			),
		)

		if activeDays > 0 {
			average := total / int64(activeDays)

			fmt.Printf(
				"Daily avg   : %s\n",
				formatDuration(
					time.Duration(average)*time.Second,
				),
			)
		}

		fmt.Println()
	},
}

func monthBar(duration time.Duration) string {
	minutes := int(duration.Minutes())

	const maxMinutes = 600

	width := 20

	filled := minutes * width / maxMinutes

	if filled > width {
		filled = width
	}

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

func init() {
	rootCmd.AddCommand(monthCmd)
}