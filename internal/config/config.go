package config

import (
	"log"
	"regexp"
	"time"

	"github.com/spf13/viper"
)

var (
	Viper     = viper.NewWithOptions(viper.KeyDelimiter("::"))
	Config    config
	CheckMode bool
)

type config struct {
	SyncFrequency              time.Duration
	ContainerTtlCreated        time.Duration
	MatchContainersLabelsRegex map[string][]*regexp.Regexp
	//ExcludeContainersRegex []string `yaml:"excludeContainersRegex"`
}

func FromViperConfig() {
	Viper.SetEnvPrefix("DC")
	Viper.AutomaticEnv()
	Viper.SetConfigType("yaml")
	Viper.SetConfigFile(Viper.GetString("config"))

	CheckMode = Viper.GetBool("check")
	if CheckMode {
		log.Println("check mode is activated")
	}

	_ = Viper.BindEnv("sync_frequency")
	Viper.SetDefault("sync_frequency", "60s")

	if err := Viper.ReadInConfig(); err == nil {
		log.Println("using configuration file:", Viper.ConfigFileUsed())
	} else {
		log.Fatal("fatal error config file: %w", err)
	}

	// Fill config struct from Viper object
	Config.SyncFrequency = Viper.GetDuration("sync_frequency")
	Config.ContainerTtlCreated = Viper.GetDuration("container_ttl_created")

	// Compile regex patterns and fill Config.MatchContainersLabelsRegex
	Config.MatchContainersLabelsRegex = make(map[string][]*regexp.Regexp)
	for key, regexList := range Viper.GetStringMapStringSlice("match_containers_labels_regex") {
		var regxs []*regexp.Regexp
		for _, regexStr := range regexList {
			regx, err := regexp.Compile(regexStr)
			if err != nil {
				log.Fatalf("failed to compile regex pattern '%s': %s", regexStr, err)
			}
			regxs = append(regxs, regx)
		}
		Config.MatchContainersLabelsRegex[key] = regxs
	}
}
