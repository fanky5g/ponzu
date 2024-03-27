package generate

import (
	"errors"
	"github.com/fanky5g/ponzu/content"
	contentGenerator "github.com/fanky5g/ponzu/content/generator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

func defineFlags(flagSet *flag.FlagSet, cwd string) {
	flagSet.String("target_package", "entities", "Target package name of content. Defaults to \"entities\"")
	flagSet.String("target_root_path", cwd, "Root path of directory to generate content. Defaults to cwd.")
	flagSet.String(
		"target_base_path",
		"entities",
		"Base path of directory relative to root path to generate content. Defaults to entities.",
	)
}

func getGeneratorArgs(flagSet *flag.FlagSet, args []string) (*GeneratorArgs, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	defineFlags(flagSet, cwd)
	if err = flagSet.Parse(args); err != nil {
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

	if err = viper.BindPFlags(flagSet); err != nil {
		return nil, err
	}

	return &GeneratorArgs{
		Arguments: args,
		Target: contentGenerator.Target{
			Path: contentGenerator.Path{
				Root: viper.GetString("target_root_path"),
				Base: viper.GetString("target_base_path"),
			},
			Package: viper.GetString("target_package"),
		},
		ContentType: content.TypeContent,
	}, nil
}

var generateCmd = &cobra.Command{
	Use:     "generate <generator type (,...fields)>",
	Aliases: []string{"gen", "g"},
	Short:   "generate boilerplate code for various Ponzu components",
	Long: `Generate boilerplate code for various Ponzu components, such as 'entities'.

The command above will generate a file 'entities/review.go' with boilerplate
methods, as well as struct definition, and corresponding field tags like:

type Review struct {
	Title  string   ` + "`json:" + `"title"` + "`" + `
	Body   string   ` + "`json:" + `"body"` + "`" + `
	Rating int      ` + "`json:" + `"rating"` + "`" + `
	Tags   []string ` + "`json:" + `"tags"` + "`" + `
}

The generate command will intelligently parse more sophisticated field names
such as 'field_name' and convert it to 'FieldName' and vice versa, only where
appropriate as per common Go idioms. Errors will be reported, but successful
generate commands return nothing.`,
	Example: `$ ponzu gen entities review title:"string" body:"string" rating:"int" tags:"[]string"`,
}

var contentCmd = &cobra.Command{
	Use:     "content <namespace> <field> <field>...",
	Aliases: []string{"c"},
	Short:   "generates a new entities type",
	RunE: func(cmd *cobra.Command, args []string) error {
		generatorArgs, err := getGeneratorArgs(cmd.Flags(), args)
		if err != nil {
			return err
		}

		g, err := NewGenerator(*generatorArgs)
		if err != nil {
			return err
		}

		return g.Generate()
	},
}

var plainTypeCmd = &cobra.Command{
	Use:     "type <namespace> <field> <field>...",
	Aliases: []string{"t"},
	Short:   "generates a new type",
	RunE: func(cmd *cobra.Command, args []string) error {
		generatorArgs, err := getGeneratorArgs(cmd.Flags(), args)
		if err != nil {
			return err
		}

		g, err := NewGenerator(*generatorArgs)
		if err != nil {
			return err
		}

		return g.Generate()
	},
}

var fieldCollectionTypeCmd = &cobra.Command{
	Use:     "field-collection <namespace> <field> <field>...",
	Aliases: []string{"fc"},
	Short:   "generates a new field-collection type",
	RunE: func(cmd *cobra.Command, args []string) error {
		generatorArgs, err := getGeneratorArgs(cmd.Flags(), args)
		if err != nil {
			return err
		}

		g, err := NewGenerator(*generatorArgs)
		if err != nil {
			return err
		}

		return g.Generate()
	},
}

func RegisterCommandRecursive(parent *cobra.Command) {
	generateCmd.AddCommand(contentCmd, plainTypeCmd, fieldCollectionTypeCmd)
	parent.AddCommand(generateCmd)
}
