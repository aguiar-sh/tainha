package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Log headers
	log.Println("Received request headers:")
	for name, values := range r.Header {
		log.Printf("  %s: %v", name, values)
	}

	// Get name and role from headers
	name := r.Header.Get("X-Name")
	role := r.Header.Get("X-Role")

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create channel for client disconnect detection
	notify := r.Context().Done()
	go func() {
		<-notify
		log.Println("Client disconnected")
	}()

	// Send events every 2 seconds
	for {
		select {
		case <-notify:
			return
		default:
			// Create event data with name and role
			event := fmt.Sprintf("data: {\"time\": \"%s\", \"message\": \"Test SSE Event\", \"name\": \"%s\", \"role\": \"%s\"}\n\n",
				time.Now().Format(time.RFC3339), name, role)

			// Write to response
			_, err := fmt.Fprint(w, event)
			if err != nil {
				log.Printf("Error writing to stream: %v", err)
				return
			}

			// Flush the response writer
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	port := flag.Int("port", 3000, "Port to run the SSE server on")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

	// Serve the client HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test/sseServer/sse_client.html")
	})

	http.HandleFunc("/sse", sseHandler)
	log.Printf("Starting SSE test server on %s\n", addr)
	log.Printf("Visit http://localhost%s to access the client\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
