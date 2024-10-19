package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aguiar-sh/tainha/internal/config"
	"github.com/aguiar-sh/tainha/internal/router"
)

func main() {
	configPath := flag.String("config", "./config/config.yaml", "Path to configuration file")
	addr := flag.String("addr", ":8080", "Address to listen on")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	r, err := router.SetupRouter(cfg)
	if err != nil {
		log.Fatalf("Erro ao configurar roteador: %v", err)
	}

	// Log all routes
	for _, route := range cfg.Routes {
		log.Printf("Route: Method=%s, Path=%s, Service=%s", route.Method, route.Path, route.Service)
		for _, mapping := range route.Mapping {
			log.Printf("  Mapping: Path=%s, Service=%s, Tag=%s", mapping.Path, mapping.Service, mapping.Tag)
		}
	}

	fmt.Printf("API Gateway escutando em %s\n", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
