package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/session"
)

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Show this week's productivity",

	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()

		weekday := int(now.Weekday())

		if weekday == 0 {
			weekday = 7
		}

		monday := now.AddDate(
			0,
			0,
			-(weekday - 1),
		)

		fmt.Println()
		fmt.Printf(
			"Week — %s to %s\n",
			monday.Format("02 Jan"),
			monday.AddDate(0, 0, 6).Format("02 Jan"),
		)

		fmt.Println("============================================")

		var total int64
		var study int64
		var work int64

		for i := 0; i < 7; i++ {
			date := monday.AddDate(0, 0, i)

			stats, err := session.GetDailyStats(date)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			total += stats.TotalSeconds
			study += stats.StudySeconds
			work += stats.WorkSeconds

			duration := time.Duration(stats.TotalSeconds) * time.Second

			fmt.Printf(
				"%s  %-8s\n",
				date.Format("Mon 02"),
				formatDuration(duration),
			)
		}

		fmt.Println("--------------------------------------------")

		fmt.Printf(
			"Study : %s\n",
			formatDuration(time.Duration(study)*time.Second),
		)

		fmt.Printf(
			"Work  : %s\n",
			formatDuration(time.Duration(work)*time.Second),
		)

		fmt.Printf(
			"Total : %s\n",
			formatDuration(time.Duration(total)*time.Second),
		)

		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(weekCmd)
}