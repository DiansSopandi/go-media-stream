package cmd

import (
	"os"

	"github.com/DiansSopandi/media_stream/bootstrap"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ems-queue",
	Short: "ems-queue is a waiting room for event management system.",
	Long:  `ems is a waiting room for event management system. It is a simple queue system that allows you to manage your events and attendees.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.ServerInitialize()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ems.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
