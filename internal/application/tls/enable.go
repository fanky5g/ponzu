// Package tls provides the functionality to Ponzu systems to encrypt HTTP traffic
// through the ability to generate self-signed certificates for local development
// and fetch/update production certificates from Let's Encrypt.
package tls

import (
	"crypto/tls"
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/util"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// newManager attempts to locate or create the cert cache directory and the
// certs for TLS encryption and returns an autocert.Manager
func (s *Service) newManager() autocert.Manager {
	cache := autocert.DirCache(filepath.Join(config.TlsDir(), "certs"))
	if _, err := os.Stat(string(cache)); os.IsNotExist(err) {
		err = os.MkdirAll(string(cache), os.ModePerm|os.ModeDir)
		if err != nil {
			log.Fatalln("Couldn't create cert directory at", cache)
		}
	}

	cfg, err := s.config.Latest()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	// get host/domain and email from GetConfig to use for TLS request to Letsencrypt.
	// we will fail fatally if either are not found since Let's Encrypt will rate-limit
	// and sending incomplete requests is wasteful and guaranteed to fail its check
	host, err := util.StringFieldByJSONTagName(cfg, "domain")
	if err != nil {
		log.Fatalf("Error identifying host/domain during TLS set-up: %v", err)
	}

	if host == "" {
		log.Fatalln("No 'domain' field set in Configuration. Please add a domain before attempting to make certificates.")
	}
	fmt.Println("Using", host, "as host/domain for certificate...")
	fmt.Println("NOTE: if the host/domain is not configured properly or is unreachable, HTTPS set-up will fail.")

	email, err := util.StringFieldByJSONTagName(cfg, "admin_email")
	if err != nil {
		log.Fatalln("Error identifying controllers email during TLS set-up.", err)
	}

	if email == "" {
		log.Fatalln("No 'admin_email' field set in Configuration. Please add an controllers email before attempting to make certificates.")
	}
	fmt.Println("Using", email, "as contact email for certificate...")

	return autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       cache,
		HostPolicy:  autocert.HostWhitelist(host),
		RenewBefore: time.Hour * 24 * 30,
		Email:       string(email),
	}
}

// Enable runs the setup for creating or locating production certificates and
// starts the TLS server
func (s *Service) Enable() {
	m := s.newManager()

	cfg, err := s.config.Latest()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	httpsPort, err := util.StringFieldByJSONTagName(cfg, "https_port")
	if err != nil {
		log.Fatalf("Failed to get https_port: %v", err)
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", httpsPort),
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}

	// launch http listener for "http-01" ACME challenge
	go func() {
		err := http.ListenAndServe(":http", m.HTTPHandler(nil))
		if err != nil {

		}
	}()

	log.Fatalln(server.ListenAndServeTLS("", ""))
}
