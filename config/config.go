package config

import (
	"errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

type Paths struct {
	PublicPath string
	DataDir    string
}

type ServeConfig struct {
	HttpsPort int
	HttpPort  int
	Bind      string
	DevHttps  bool
	Https     bool
}

type Config struct {
	Paths          Paths
	ServeConfig    ServeConfig
	DatabaseDriver string
}

func defineFlags(flagSet *flag.FlagSet, workingDir string) {
	flagSet.String("public_path", "", "Optional URL path to serve ponzu on")
	flagSet.String("bind", "localhost", "address for ponzu to bind the HTTP(S) server")
	flagSet.Int("https_port", 443, "port for ponzu to bind its HTTPS listener")
	flagSet.Int("port", 8080, "port for ponzu to bind its HTTP listener")
	flagSet.Bool("https", false, "enable automatic TLS/SSL certificate management")
	flagSet.Bool("dev_https", false, "[dev environment] enable automatic TLS/SSL certificate management")
	flagSet.String(
		"data_dir",
		workingDir,
		"directory where application data should be stored. Defaults to working directory",
	)
	flagSet.String("database_driver", "", "Database driver to use.")
}

func Get() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	flags := flag.NewFlagSet("config", flag.ExitOnError)
	defineFlags(flags, cwd)

	if err = flags.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	viper.SetConfigName("ponzu")
	viper.SetConfigType("props")
	viper.AddConfigPath(cwd)
	err = viper.ReadInConfig()
	if err != nil && errors.As(err, &viper.ConfigFileNotFoundError{}) {
		log.Info("config file not found. will default to provided flags")
		err = nil
	}

	if err = viper.BindPFlags(flags); err != nil {
		return nil, err
	}

	return &Config{
		Paths: Paths{
			PublicPath: viper.GetString("public_path"),
		},
		ServeConfig: ServeConfig{
			HttpsPort: viper.GetInt("https_port"),
			HttpPort:  viper.GetInt("port"),
			Bind:      viper.GetString("bind"),
			DevHttps:  viper.GetBool("dev_https"),
			Https:     viper.GetBool("https"),
		},
		DatabaseDriver: viper.GetString("database_driver"),
	}, nil
}

func New() (*Config, error) {
	return Get()
}
