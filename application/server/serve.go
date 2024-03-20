package server

import (
	"fmt"
	"net/http"
)

func (server *server) Serve() error {
	// cannot run production HTTPS and development HTTPS together
	if server.cfg.DevHttps {
		fmt.Println("Enabling self-signed HTTPS... [DEV]")

		go server.tlsService.EnableDev()
		fmt.Println("Server listening on https://localhost:10443 for requests... [DEV]")
		fmt.Println("----")
		fmt.Println("If your browser rejects HTTPS requests, try allowing insecure connections on localhost.")
		fmt.Println("on Chrome, visit chrome://flags/#allow-insecure-localhost")

	} else if server.cfg.Https {
		fmt.Println("Enabling HTTPS...")

		go server.tlsService.Enable()
		fmt.Printf(
			"Server listening on :%s for HTTPS requests...\n",
			server.configCache.GetByKey("https_port").(string),
		)
	}

	fmt.Printf("Server listening at %s:%d for HTTP requests...\n", server.cfg.Bind, server.cfg.HttpPort)
	fmt.Println("\nVisit '/' to get started.")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", server.cfg.Bind, server.cfg.HttpPort), server.mux)
}
