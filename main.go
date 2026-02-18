//
// rubrik-exporter
//
// Exports metrics from rubrik backup for prometheus
//
// License: Apache License Version 2.0,
// Organization: Claranet GmbH
// Author: Martin Weber <martin.weber@de.clara.net>
//

package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/Gattancha-Computer-Services/rubrik-exporter/rubrik"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var rubrikAPI *rubrik.Rubrik
var vmIDNameMap map[string]string

var (
	namespace                    = "rubrik"
	rubrikURL                    = flag.String("rubrik.url", "", "Rubrik URL to connect https://rubrik.local.host")
	rubrikUser                   = flag.String("rubrik.username", "", "Rubrik API User")
	rubrikPassword               = flag.String("rubrik.password", "", "Rubrik API User Password")
	rubrikServiceAccountClientID = flag.String("rubrik.service-account-client-id", "", "Rubrik Service Account Client ID")
	rubrikServiceAccountClientSecret = flag.String("rubrik.service-account-client-secret", "", "Rubrik Service Account Client Secret")
	listenAddress                = flag.String("listen-address", ":9477", "The address to listen on for HTTP requests.")
)

func main() {
	flag.Parse()

	log.Print("Create Rubrik Exporter instance")
	rubrikAPI = rubrik.NewRubrik(*rubrikURL, *rubrikUser, *rubrikPassword, *rubrikServiceAccountClientID, *rubrikServiceAccountClientSecret)

	prometheus.MustRegister(NewRubrikStatsExport())
	prometheus.MustRegister(NewVMStatsExport())
	prometheus.MustRegister(NewArchiveLocation())
	prometheus.MustRegister(NewManagedVolume())

	metricsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Metrics request from %s - User-Agent: %s", r.RemoteAddr, r.Header.Get("User-Agent"))
		promhttp.Handler().ServeHTTP(w, r)
	})

	// Serve metrics at both /metrics and / for compatibility
	http.Handle("/metrics", metricsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If request is for metrics (Accept header or direct access), serve metrics
		acceptHeader := r.Header.Get("Accept")
		userAgent := r.Header.Get("User-Agent")
		
		// Serve metrics if it's from Prometheus/Grafana or has json in Accept header
		if strings.Contains(acceptHeader, "text/plain") || 
		   strings.Contains(acceptHeader, "application/json") ||
		   strings.Contains(userAgent, "Prometheus") ||
		   strings.Contains(userAgent, "Grafana") ||
		   r.URL.Path == "/" {
			metricsHandler.ServeHTTP(w, r)
			return
		}
		
		// Otherwise show HTML welcome page for browsers
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<html><head><title>Rubrik Exporter</title></head><body><h1>Rubrik Exporter</h1><p><a href="/metrics">Metrics</a></p></body></html>`))
	})

	log.Printf("Starting Server: %s", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
