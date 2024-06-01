package generate

import (
	"github.com/fanky5g/ponzu/generator"
	"github.com/spf13/cobra"
)

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
	Example: `$ ponzu gen content review title:"string" body:"string" rating:"int" tags:"[]string"`,
}

var contentCmd = &cobra.Command{
	Use:     "content <namespace> <field> <field>...",
	Aliases: []string{"c"},
	Short:   "generates a new content type",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(cmd.Flags(), args); err != nil {
			return err
		}

		return RunGenerator(generator.Content, args)
	},
}

var plainTypeCmd = &cobra.Command{
	Use:     "type <namespace> <field> <field>...",
	Aliases: []string{"t"},
	Short:   "generates a new type",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(cmd.Flags(), args); err != nil {
			return err
		}

		return RunGenerator(generator.Plain, args)
	},
}

var fieldCollectionTypeCmd = &cobra.Command{
	Use:     "field-collection <namespace> <field> <field>...",
	Aliases: []string{"fc"},
	Short:   "generates a new field-collection type",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(cmd.Flags(), args); err != nil {
			return err
		}

		return RunGenerator(generator.FieldCollection, args)
	},
}

func RegisterCommandRecursive(parent *cobra.Command) {
	generateCmd.AddCommand(contentCmd, plainTypeCmd, fieldCollectionTypeCmd)
	parent.AddCommand(generateCmd)
}
