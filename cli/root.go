package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand *cobra.Command = &cobra.Command{
	Use:   "forgequeue",
	Short: "forgequeue helps",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("root command")
		return nil
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCommand.AddCommand(workerCommand)
	rootCommand.AddCommand(schedulerCommand)
	rootCommand.AddCommand(recovererCommand)
}
