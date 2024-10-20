package router

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/aguiar-sh/tainha/internal/config"
	"github.com/aguiar-sh/tainha/internal/mapper"
	"github.com/aguiar-sh/tainha/internal/proxy"
	"github.com/aguiar-sh/tainha/internal/util"
	"github.com/gorilla/mux"
)

func SetupRouter(cfg *config.Config) (*mux.Router, error) {
	r := mux.NewRouter()

	for _, route := range cfg.Routes {

		path, protocol := util.PathProtocol(route.Service)
		servicePath := fmt.Sprintf("%s://%s", protocol, path)

		reverseProxy, err := proxy.NewReverseProxy(servicePath)
		if err != nil {
			log.Fatalf("Erro ao criar proxy para %s: %v", route.Path, err)
		}

		fullPath := fmt.Sprintf("%s%s", cfg.BaseConfig.BasePath, route.Route)

		r.HandleFunc(fullPath, func(w http.ResponseWriter, req *http.Request) {
			log.Println("Request received for:", req.URL.Path)

			// Extract path parameters using the utility function
			params := util.ExtractPathParams(route.Path)
			vars := mux.Vars(req)

			// Construct the target path dynamically
			targetPath := route.Path
			for _, param := range params {
				value, ok := vars[param]
				if !ok {
					http.Error(w, fmt.Sprintf("Parameter %s not found in request path", param), http.StatusBadRequest)
					return
				}
				// Replace the placeholder with the actual value
				targetPath = strings.Replace(targetPath, fmt.Sprintf("{%s}", param), value, -1)
			}

			// Parse the target path to separate path and query
			parsedURL, err := url.Parse(targetPath)
			if err != nil {
				log.Printf("Error parsing target path: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			req.URL.Path = parsedURL.Path
			req.URL.RawQuery = parsedURL.RawQuery
			req.URL.Host = route.Service
			req.URL.Scheme = protocol

			// Capture the response from the reverse proxy
			rec := httptest.NewRecorder()
			reverseProxy.ServeHTTP(rec, req)

			// Read the response body
			respBody := rec.Body.Bytes()

			if rec.Code < 200 || rec.Code >= 300 {
				for k, v := range rec.Header() {
					w.Header()[k] = v
				}
				w.WriteHeader(rec.Code)
				w.Write(respBody)
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
	}

	return r, nil
}
