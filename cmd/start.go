package cmd

import (
	"log"
	"path/filepath"
	"time"

	"github.com/DiansSopandi/media_stream/bootstrap"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/spf13/cobra"
)

var cfgFile string
var tz string
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Long:  `Start the application with the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		// fmt.Println("🚀 App started from startCmd")
		// Panggil aplikasi atau server kamu di sini, misalnya:
		bootstrap.ServerInitialize()
	},
}

// init is called before the command is executed.
// @return void
func init() {
	rootCmd.AddCommand(startCmd)
	cobra.OnInitialize(initConfig)

	startCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "/path/to/config/env.conf")
	startCmd.PersistentFlags().StringVar(&tz, "timezone", "", "Asia/Jakarta")
}

func loadConfig() {
	// Load config file
	if cfgFile != "" {
		cfgPath := filepath.Dir(cfgFile)
		cfgFilename := filepath.Base(cfgFile)

		pkg.LoadConfig(cfgFilename, cfgPath)
	} else {
		pkg.LoadConfig("env.conf", "./")
	}

	// Set timezone
	if tz != "" {
		pkg.Cfg.Application.Timezone = tz
	}

	loc, err := time.LoadLocation(pkg.Cfg.Application.Timezone)
	if err != nil {
		log.Fatalf("Error loading timezone, %v", err)
	}
	time.Local = loc
}

// initConfig initializes the configuration.
// @return void
func initConfig() {
	// Load config file
	loadConfig()

	// Connect to database
	// database := db.InitDatabase()
	// defer database.Close()

	// pkg.DatabaseConnect()
	// pkg.ConnectMaticore()
	// pkg.Pooling()
}
