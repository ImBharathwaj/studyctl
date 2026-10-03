package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/database"
)

var (
	logDuration string
	logType     string
	logNote     string
)

var logCmd = &cobra.Command{
	Use:   "log [task]",
	Short: "Manually log a completed session",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		duration, err := parseDuration(logDuration)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if logType != "study" && logType != "work" {
			fmt.Println("Error: type must be study or work")
			return
		}

		end := time.Now()
		start := end.Add(-duration)

		_, err = database.DB.Exec(`
			INSERT INTO sessions
			(task, type, start_time, end_time,
			 duration_seconds, notes)
			VALUES (?, ?, ?, ?, ?, ?)
		`,
			args[0],
			logType,
			start,
			end,
			int64(duration.Seconds()),
			logNote,
		)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Session logged")
		fmt.Println("----------------------------")
		fmt.Printf("Task     : %s\n", args[0])
		fmt.Printf("Type     : %s\n", logType)
		fmt.Printf("Duration : %s\n", formatDuration(duration))
	},
}

func init() {
	logCmd.Flags().StringVar(
		&logDuration,
		"duration",
		"",
		"Duration, e.g. 90m or 2h",
	)

	logCmd.MarkFlagRequired("duration")

	logCmd.Flags().StringVar(
		&logType,
		"type",
		"study",
		"study or work",
	)

	logCmd.Flags().StringVar(
		&logNote,
		"note",
		"",
		"Session note",
	)

	rootCmd.AddCommand(logCmd)
}