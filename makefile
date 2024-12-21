run:
	go run cmd/gateway/main.go

# Run SSE test server
sse-server:
	@echo "Starting SSE test server on :3000/sse"
	@go run test/sseServer/sse_server.go