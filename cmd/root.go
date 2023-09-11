package cmd

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"docker-cleaner/internal/config"
	"docker-cleaner/internal/docker"

	"github.com/spf13/cobra"
)

// Build information. Populated at build-time.
var (
	program   string
	version   string
	commit    string
	buildDate string
	goVersion = runtime.Version()
	goOS      = runtime.GOOS
	goArch    = runtime.GOARCH
)

// Version format
var versionFormat = `version %s
  git commit: %s
  build date: %s
  go version: %s
  platform:   %s`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     program,
	Short:   "Remove abandoned docker-compose projects.",
	Version: fmt.Sprintf(versionFormat, version, commit, buildDate, goVersion, goOS+"/"+goArch),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Docker cleaner run with sync_frequency:", config.Config.SyncFrequency)
		cleanerTicker := time.NewTicker(config.Config.SyncFrequency)
		for {
			select {
			case <-cleanerTicker.C:
				docker.Cleaner()
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")

	flags := rootCmd.Flags()
	flags.StringP("config", "c", program+".yml", "path to config file")
	flags.BoolP("check", "C", false, "don't make any changes; instead, try to predict some of the changes that may occur")

	// Bind all cmd flags to viper for gather env vars
	_ = config.Viper.BindPFlags(rootCmd.Flags())

	cobra.OnInitialize(config.FromViperConfig)
}
