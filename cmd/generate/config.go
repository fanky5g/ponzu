package generate

import (
	"errors"
	"github.com/fanky5g/ponzu/generator"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

func defineFlags(flagSet *flag.FlagSet, cwd string) {
	// Generated content files
	flagSet.String("target_package", "entities", "Target package name of content. Defaults to \"entities\"")
	flagSet.String("target_root_path", cwd, "Root path of directory to generate content. Defaults to cwd.")
	flagSet.String(
		"target_base_path",
		"entities",
		"Base path of directory relative to root path to generate content. Defaults to entities.",
	)

	// Generated model files
	flagSet.String("models_target_package", "models", "Target package name of models. Defaults to \"models\"")
	flagSet.String("models_target_root_path", cwd, "Root path of directory to generate models. Defaults to cwd.")
	flagSet.String(
		"models_target_base_path",
		"models",
		"Base path of directory relative to root path to generate models. Defaults to models.",
	)
}

func initConfig(flagSet *flag.FlagSet, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	defineFlags(flagSet, cwd)
	if err = flagSet.Parse(args); err != nil {
		return err
	}

	viper.SetConfigName("ponzu")
	viper.SetConfigType("props")
	viper.AddConfigPath(cwd)
	err = viper.ReadInConfig()
	if err != nil && errors.As(err, &viper.ConfigFileNotFoundError{}) {
		log.Info("config file not found. will default to provided flags")
		err = nil
	}

	if err = viper.BindPFlags(flagSet); err != nil {
		return err
	}

	return nil
}

func getContentConfig(contentType generator.Type) generator.Config {
	return generator.Config{
		Target: generator.Target{
			Path: generator.Path{
				Root: viper.GetString("target_root_path"),
				Base: viper.GetString("target_base_path"),
			},
			Package: viper.GetString("target_package"),
		},
		Type: contentType,
	}
}

func getModelConfig() generator.Config {
	return generator.Config{
		Target: generator.Target{
			Path: generator.Path{
				Root: viper.GetString("models_target_root_path"),
				Base: viper.GetString("models_target_base_path"),
			},
			Package: viper.GetString("models_target_package"),
		},
	}
}
