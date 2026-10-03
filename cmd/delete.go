package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yourname/studyctl/internal/session"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a session",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)

		if err != nil {
			fmt.Println("Invalid session ID")
			return
		}

		err = session.Delete(id)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Session #%d deleted.\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}