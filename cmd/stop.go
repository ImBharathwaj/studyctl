package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/session"
)

var note string

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current session",

	Run: func(cmd *cobra.Command, args []string) {
		s, err := session.Stop(note)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		duration := time.Duration(s.DurationSeconds) * time.Second

		fmt.Println()
		fmt.Println("Session completed")
		fmt.Println("----------------------------")
		fmt.Printf("Task     : %s\n", s.Task)
		fmt.Printf("Type     : %s\n", s.Type)
		fmt.Printf("Duration : %s\n", formatDuration(duration))

		if note != "" {
			fmt.Printf("Note     : %s\n", note)
		}

		fmt.Println("----------------------------")
	},
}

func init() {
	stopCmd.Flags().StringVar(
		&note,
		"note",
		"",
		"Add a note to the session",
	)

	rootCmd.AddCommand(stopCmd)
}