package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/database"
	"github.com/yourname/studyctl/internal/session"
)

var goalSetCmd = &cobra.Command{
	Use:   "set [duration]",
	Short: "Set today's target",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		duration, err := parseDuration(args[0])

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = database.SetGoal(time.Now(), int64(duration.Seconds()))

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf(
			"Daily goal set to %s\n",
			formatDuration(duration),
		)
	},
}

var goalShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show today's goal",

	Run: func(cmd *cobra.Command, args []string) {
		goal, err := database.GetGoal(time.Now())

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if goal == 0 {
			fmt.Println("No daily goal configured.")
			return
		}

		stats, err := session.GetDailyStats(time.Now())

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		target := time.Duration(goal) * time.Second
		current := time.Duration(stats.TotalSeconds) * time.Second

		percentage := float64(current) / float64(target) * 100

		fmt.Println()
		fmt.Println("Today's goal")
		fmt.Println("--------------------------------------------")
		fmt.Printf("Target   : %s\n", formatDuration(target))
		fmt.Printf("Current  : %s\n", formatDuration(current))
		fmt.Printf("Progress : %.0f%%\n", percentage)
		fmt.Printf("[%s]\n", progressBar(percentage, 30))

		if current >= target {
			fmt.Println("Goal completed.")
		} else {
			remaining := target - current
			fmt.Printf(
				"Remaining: %s\n",
				formatDuration(remaining),
			)
		}
	},
}

var goalCmd = &cobra.Command{
	Use:   "goal",
	Short: "Manage daily goals",
}

func parseDuration(value string) (time.Duration, error) {
	if strings.HasSuffix(value, "h") {
		raw := strings.TrimSuffix(value, "h")

		hours, err := strconv.ParseFloat(raw, 64)

		if err != nil {
			return 0, err
		}

		return time.Duration(hours * float64(time.Hour)), nil
	}

	if strings.HasSuffix(value, "m") {
		raw := strings.TrimSuffix(value, "m")

		minutes, err := strconv.Atoi(raw)

		if err != nil {
			return 0, err
		}

		return time.Duration(minutes) * time.Minute, nil
	}

	return 0, fmt.Errorf(
		"duration must look like 8h or 90m",
	)
}

func init() {
	goalCmd.AddCommand(goalSetCmd)
	goalCmd.AddCommand(goalShowCmd)

	rootCmd.AddCommand(goalCmd)
}