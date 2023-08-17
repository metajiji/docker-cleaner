package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Path to config file
var configFile string

// Build information. Populated at build-time.
var (
	Program   string
	Version   string
	Commit    string
	BuildDate string
	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
)
var versionFormat = `version %s
  git commit: %s
  build date: %s
  go version: %s
  platform:   %s`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     Program,
	Short:   "Remove abandoned docker-compose projects.",
	Version: fmt.Sprintf(versionFormat, Version, Commit, BuildDate, GoVersion, GoOS+"/"+GoArch),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: Program started...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", Program+".yml", "path to config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("DC")
	handleError(viper.BindEnv("API_KEY"))
	handleError(viper.BindEnv("API_SECRET"))
	handleError(viper.BindEnv("USERNAME"))
	handleError(viper.BindEnv("PASSWORD"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using configuration file: ", viper.ConfigFileUsed())
	}
}

// TODO Remove this
func handleError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
