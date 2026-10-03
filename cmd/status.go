package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/session"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current session",

	Run: func(cmd *cobra.Command, args []string) {

		s, err := session.Active()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if s == nil {
			fmt.Println("No active session.")
			return
		}

		elapsed := time.Since(s.StartTime)

		fmt.Println("Current session")
		fmt.Println("----------------------------")
		fmt.Printf("Task      : %s\n", s.Task)
		fmt.Printf("Started   : %s\n",
			s.StartTime.Format("15:04:05"))
		fmt.Printf("Elapsed   : %s\n",
			elapsed.Round(time.Second))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}