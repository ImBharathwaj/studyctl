package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/session"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show today's sessions",

	Run: func(cmd *cobra.Command, args []string) {
		rows, err := session.GetSessions(time.Now())

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		defer rows.Close()

		fmt.Println()
		fmt.Println("Today's sessions")
		fmt.Println("====================================================")

		for rows.Next() {
			var (
				id       int64
				task     string
				kind     string
				start    time.Time
				end      time.Time
				duration int64
			)

			err := rows.Scan(
				&id,
				&task,
				&kind,
				&start,
				&end,
				&duration,
			)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Printf(
				"#%-3d %-28s %-6s %s → %s  %s\n",
				id,
				task,
				kind,
				start.Format("15:04"),
				end.Format("15:04"),
				formatDuration(
					time.Duration(duration)*time.Second,
				),
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}