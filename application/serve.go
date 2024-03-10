package application

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/internal/services/tls"
	"log"
	"net/http"
)

func (app *application) Serve() error {
	cfg := app.serveConfig
	if cfg == nil {
		return errors.New("invalid serve config")
	}

	err := app.repositories.Config.PutConfig("https_port", fmt.Sprintf("%d", cfg.HttpsPort))
	if err != nil {
		log.Fatalln("System failed to save config. Please try to run again.", err)
	}

	// save the https port the system is listening on so internal system can make
	// HTTP api calls while in dev or production w/o adding more cli flags
	err = app.repositories.Config.PutConfig("http_port", fmt.Sprintf("%d", cfg.HttpPort))
	if err != nil {
		log.Fatalln("System failed to save config. Please try to run again.", err)
	}

	bind := cfg.Bind
	// save the bound address the system is listening on so internal system can make
	// HTTP api calls while in dev or production w/o adding more cli flags
	if bind == "" {
		bind = "localhost"
	}
	err = app.repositories.Config.PutConfig("bind_addr", bind)
	if err != nil {
		log.Fatalln("System failed to save config. Please try to run again.", err)
	}

	tlsService := app.services.Get(tls.ServiceToken).(tls.Service)

	// cannot run production HTTPS and development HTTPS together
	if cfg.DevHttps {
		fmt.Println("Enabling self-signed HTTPS... [DEV]")

		go tlsService.EnableDev()
		fmt.Println("Server listening on https://localhost:10443 for requests... [DEV]")
		fmt.Println("----")
		fmt.Println("If your browser rejects HTTPS requests, try allowing insecure connections on localhost.")
		fmt.Println("on Chrome, visit chrome://flags/#allow-insecure-localhost")

	} else if cfg.Https {
		fmt.Println("Enabling HTTPS...")

		go tlsService.Enable()
		fmt.Printf(
			"Server listening on :%s for HTTPS requests...\n",
			app.repositories.Config.Cache().GetByKey("https_port").(string),
		)
	}

	fmt.Printf("Server listening at %s:%d for HTTP requests...\n", bind, cfg.HttpPort)
	fmt.Println("\nVisit '/' to get started.")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", bind, cfg.HttpPort), nil)
}

func (app *application) ServeMux() *http.ServeMux {
	return app.mux
}
