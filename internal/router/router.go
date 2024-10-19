package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/aguiar-sh/tainha/internal/config"
	"github.com/aguiar-sh/tainha/internal/mapper"
	"github.com/aguiar-sh/tainha/internal/proxy"
	"github.com/gorilla/mux"
)

func extractPathParams(path string) []string {
	re := regexp.MustCompile(`{([^}]+)}`)
	matches := re.FindAllStringSubmatch(path, -1)
	params := make([]string, len(matches))
	for i, match := range matches {
		params[i] = match[1]
	}
	return params
}

func SetupRouter(cfg *config.Config) (*mux.Router, error) {
	r := mux.NewRouter()

	for _, route := range cfg.Routes {
		reverseProxy, err := proxy.NewReverseProxy(route.Service)
		if err != nil {
			log.Fatalf("Erro ao criar proxy para %s: %v", route.Path, err)
		}

		r.HandleFunc(route.Path, func(w http.ResponseWriter, req *http.Request) {
			log.Println("Request received for:", req.URL.Path)

			// Capture the response from the reverse proxy
			rec := httptest.NewRecorder()
			reverseProxy.ServeHTTP(rec, req)

			// Read the response body
			respBody, err := io.ReadAll(rec.Body)
			if err != nil {
				http.Error(w, "Failed to read response body", http.StatusInternalServerError)
				return
			}

			// Check if the response body is empty
			if len(respBody) == 0 {
				log.Println("Warning: Response body is empty")
			}

			response, err := mapper.Map(route, respBody)
			if err != nil {
				log.Printf("Error mapping response: %v", err)
				http.Error(w, "Failed to map response", http.StatusInternalServerError)
				return
			}

			// Copy headers from the recorder to the response writer
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}

			// Set the correct Content-Length
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(response)))

			// Write the status code
			w.WriteHeader(rec.Code)

			// Write the final response
			n, err := w.Write(response)
			if err != nil {
				log.Printf("Error writing response: %v", err)
				return
			} else if n != len(response) {
				log.Printf("Warning: not all bytes were written. Expected %d, wrote %d", len(response), n)
			}
		}).Methods(route.Method)

		// Implementar mapeamentos
		for _, m := range route.Mapping {
			mappedPath := m.Path
			mappedProxy, err := proxy.NewReverseProxy(m.Service)
			if err != nil {
				log.Fatalf("Erro ao criar proxy para mapeamento %s: %v", m.Path, err)
			}

			r.HandleFunc(mappedPath, func(w http.ResponseWriter, req *http.Request) {
				fmt.Println("Request received for:", req.URL.Path)
				// Opcional: Adicionar lógica com base na tag
				mappedProxy.ServeHTTP(w, req)
			}).Methods(route.Method)
		}
	}

	return r, nil
}
