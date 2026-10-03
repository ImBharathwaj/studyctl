package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/session"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's productivity",

	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()

		stats, err := session.GetDailyStats(now)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println()
		fmt.Printf("Today — %s\n", now.Format("02 Jan 2006"))
		fmt.Println("============================================")

		fmt.Printf(
			"Study : %s\n",
			formatDuration(time.Duration(stats.StudySeconds)*time.Second),
		)

		fmt.Printf(
			"Work  : %s\n",
			formatDuration(time.Duration(stats.WorkSeconds)*time.Second),
		)

		fmt.Printf(
			"Total : %s\n",
			formatDuration(time.Duration(stats.TotalSeconds)*time.Second),
		)

		fmt.Println()

		printHourlyTimeline(now)
		printTaskBreakdown(now)

		fmt.Println()
	},
}

func printHourlyTimeline(date time.Time) {
	hours, err := session.GetHourlyStats(date)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Hourly timeline")
	fmt.Println("--------------------------------------------")

	for _, h := range hours {
		study := time.Duration(h.StudySeconds) * time.Second
		work := time.Duration(h.WorkSeconds) * time.Second

		fmt.Printf(
			"%02d:00  Study %-7s  Work %-7s\n",
			h.Hour,
			formatDuration(study),
			formatDuration(work),
		)
	}
}

func printTaskBreakdown(date time.Time) {
	tasks, err := session.GetTaskStats(date)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println()
	fmt.Println("Task breakdown")
	fmt.Println("--------------------------------------------")

	for _, task := range tasks {
		duration := time.Duration(task.Seconds) * time.Second

		fmt.Printf(
			"%-30s %-6s %s\n",
			task.Task,
			task.Type,
			formatDuration(duration),
		)
	}
}

func init() {
	rootCmd.AddCommand(todayCmd)
}