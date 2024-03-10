package serve

import (
	"github.com/fanky5g/ponzu/application"
	"github.com/fanky5g/ponzu/database"
	"github.com/spf13/cobra"
)

var (
	bind      string
	httpsPort int
	port      int
	https     bool
	devHttps  bool
)

func defineFlags() {
	serveCmd.Flags().StringVar(&bind, "bind", "localhost", "address for ponzu to bind the HTTP(S) server")
	serveCmd.Flags().IntVar(&httpsPort, "https-port", 443, "port for ponzu to bind its HTTPS listener")
	serveCmd.Flags().IntVar(&port, "port", 8080, "port for ponzu to bind its HTTP listener")
	serveCmd.Flags().BoolVar(&https, "https", false, "enable automatic TLS/SSL certificate management")
	serveCmd.Flags().BoolVar(&devHttps, "dev-https", false, "[dev environment] enable automatic TLS/SSL certificate management")
}

var serveCmd = &cobra.Command{
	Use:     "serve [flags]",
	Aliases: []string{"s"},
	Short:   "run the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := database.New()
		if err != nil {
			panic(err)
		}

		serveConfig := &application.ServeConfig{
			HttpsPort: httpsPort,
			HttpPort:  port,
			Bind:      bind,
			DevHttps:  devHttps,
			Https:     https,
		}

		app, err := application.New(application.Config{
			ServeConfig: serveConfig,
			Database:    db,
		})

		return app.Serve()
	},
}

func RegisterCommandRecursive(parent *cobra.Command) {
	defineFlags()
	parent.AddCommand(serveCmd)
}
