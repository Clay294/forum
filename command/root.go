package command

import (
	"fmt"

	"github.com/Clay294/forum/command/start"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		log.Error().Msgf("execution of the program startup command failed:%s", err)
		return fmt.Errorf("execution of the program startup command failed")
	}

	return nil
}

var rootCmd = cobra.Command{
	Use:   "start [-f --config-file] [-l --log-file]",
	Short: "forum project api",
	Long:  "forum project api v2",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(&start.StartCmd)
}
