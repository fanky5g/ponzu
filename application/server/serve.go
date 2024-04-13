package server

import (
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"net/http"
)

func (server *server) Serve() error {
	appConfig, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to ponzu config: %v", err)
	}

	cfg, err := server.configService.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}

	// cannot run production HTTPS and development HTTPS together
	if appConfig.ServeConfig.DevHttps {
		fmt.Println("Enabling self-signed HTTPS... [DEV]")

		go server.tlsService.EnableDev()
		fmt.Println("Server listening on https://localhost:10443 for requests... [DEV]")
		fmt.Println("----")
		fmt.Println("If your browser rejects HTTPS requests, try allowing insecure connections on localhost.")
		fmt.Println("on Chrome, visit chrome://flags/#allow-insecure-localhost")

	} else if appConfig.ServeConfig.Https {
		fmt.Println("Enabling HTTPS...")

		go server.tlsService.Enable()
		fmt.Printf(
			"Server listening on :%s for HTTPS requests...\n",
			cfg.HTTPSPort,
		)
	}

	// start analytics recorder
	go server.analyticsService.StartRecorder()

	fmt.Printf("Server listening at %s:%d for HTTP requests...\n", cfg.BindAddress, cfg.HTTPPort)
	fmt.Println("\nVisit '/' to get started.")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.HTTPPort), server.mux)
}
