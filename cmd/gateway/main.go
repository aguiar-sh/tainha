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

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	r, err := router.SetupRouter(cfg)
	if err != nil {
		log.Fatalf("Error setting up router: %v", err)
	}

	port := fmt.Sprintf(":%d", cfg.BaseConfig.Port)
	addr := flag.String("addr", port, "Address to listen on")
	flag.Parse()

	fmt.Print("\n")

	// Log all routes
	for _, route := range cfg.Routes {

		fullPath := fmt.Sprintf("%s%s", cfg.BaseConfig.BasePath, route.Path)

		log.Println("\033[1;32mRoutes:\033[0m")

		log.Printf("\033[1;34m| %s - %s -> %s%s\033[0m", route.Method, fullPath, route.Service, route.Path)
		log.Println("\033[1;33m | Mappings: \033[0m")

		for _, mapping := range route.Mapping {
			log.Printf("\033[1;33m | %s - %s%s\033[0m\n\n", mapping.Tag, mapping.Service, mapping.Path)
		}
	}

	log.Printf("API Gateway listening on %s\n\n", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
