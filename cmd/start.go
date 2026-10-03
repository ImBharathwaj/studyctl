package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/session"
)

var sessionType string

var startCmd = &cobra.Command{
	Use:   "start <task>",
	Short: "Start a study or work session",

	Long: `Start a new study or work session.

Only one session can be active at a time.
Use "studyctl stop" to finish the current session.`,

	Args: cobra.ExactArgs(1),

	Example: `  studyctl start "Machine Learning"
  studyctl start "AWS MLA" --type study
  studyctl start "Client Project" --type work`,

	Run: func(cmd *cobra.Command, args []string) {
		task := args[0]

		err := session.Start(task, sessionType)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println()
		fmt.Println("Session started")
		fmt.Println("----------------------------")
		fmt.Printf("Task : %s\n", task)
		fmt.Printf("Type : %s\n", sessionType)
		fmt.Println("----------------------------")
	},
}

func init() {
	startCmd.Flags().StringVarP(
		&sessionType,
		"type",
		"t",
		"study",
		"Session type: study or work",
	)

	rootCmd.AddCommand(startCmd)
}