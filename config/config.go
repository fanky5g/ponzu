package config

import (
	"errors"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	once sync.Once
	cfg  *Config
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
	SearchDriver   string
	StorageDriver  string
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
	flagSet.String("search_driver", "", "Search driver to use.")
	flagSet.String("database_driver", "", "Database driver to use.")
	flagSet.String("storage_driver", "", "Upload file storage driver to use.")
}

func Get() (*Config, error) {
	var err error

	once.Do(func() {
		var cwd string
		cwd, err = os.Getwd()
		if err != nil {
			return
		}

		flags := flag.NewFlagSet("config", flag.ExitOnError)
		defineFlags(flags, cwd)

		if err = flags.Parse(os.Args[1:]); err != nil {
			return
		}

		viper.SetConfigName("ponzu")
		viper.SetConfigType("props")
		viper.AddConfigPath(cwd)
		err = viper.ReadInConfig()
		if err != nil {
			if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
				return
			}

			log.Info("config file not found. will default to provided flags")
			err = nil
		}

		if err = viper.BindPFlags(flags); err != nil {
			return
		}

		cfg = &Config{
			Paths: Paths{
				PublicPath: viper.GetString("public_path"),
				DataDir:    viper.GetString("data_dir"),
			},
			ServeConfig: ServeConfig{
				HttpsPort: viper.GetInt("https_port"),
				HttpPort:  viper.GetInt("port"),
				Bind:      viper.GetString("bind"),
				DevHttps:  viper.GetBool("dev_https"),
				Https:     viper.GetBool("https"),
			},
			DatabaseDriver: viper.GetString("database_driver"),
			SearchDriver:   viper.GetString("search_driver"),
			StorageDriver:  viper.GetString("storage_driver"),
		} 
	})

	return cfg, err
}

func New() (*Config, error) {
	return Get()
}
