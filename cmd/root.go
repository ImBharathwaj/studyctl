package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/database"
)

var rootCmd = &cobra.Command{
	Use:   "studyctl",
	Short: "A Linux CLI for tracking study and work time",
	Long: `studyctl is a local-first productivity tracker
for recording study and work sessions.

Track your time, monitor daily goals,
and analyze your productivity over days,
weeks, and months.`,

	Example: `  studyctl start "Machine Learning"
  studyctl stop
  studyctl today
  studyctl week
  studyctl month
  studyctl goal set 8h
  studyctl goal show`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := database.Init(); err != nil {
			fmt.Fprintln(os.Stderr, "Database error:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}