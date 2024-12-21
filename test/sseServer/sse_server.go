package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
			// Create event data
			event := fmt.Sprintf("data: {\"time\": \"%s\", \"message\": \"Test SSE Event\"}\n\n",
				time.Now().Format(time.RFC3339))

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
	http.HandleFunc("/sse", sseHandler)
	log.Println("Starting SSE test server on :3000/sse")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
