![Tainha logo](https://github.com/user-attachments/assets/a1286b71-5b0b-4d1e-90c6-177dd9ca5fe5)

**Tainha** is an open-source API Gateway developed in Go (Golang), inspired by the rich culture of Florianópolis. Designed to be an affordable and efficient solution, Tainha simplifies routing HTTP requests to various backend services with ease and flexibility.

## 🌟 Features

- [x] **Simple Routing:** Direct requests to backend services based on defined routes.
- [x] **Request with Parameters:** Handle routes with dynamic parameters like `{userId}`.
- [x] **Request with Queries:** Support query parameters in routes and mappings.
- [x] **Circular Requesting:** Allow chained API calls between services.
- [ ] **Mapping Cache:** Implement caching for mapped responses to improve performance.
- [ ] **Rate Limiting:** Control the rate of incoming requests to protect backend services.
- [ ] **JWT Validation with JWK:** Authenticate requests using JSON Web Tokens and JSON Web Keys.
- [x] **Open Source:** Contribute and customize Tainha to fit your needs.

We are on developing, feel free to collaborate

## 🚀 How to test

To test the application, you can use the *json-server* with `test/db.json` file.


```bash
npx json-server --watch test/db.json --port 3000
```

Then, you can start the Tainha application with the following command:

```bash
make run
```

or

```bash
go run cmd/gateway/main.go
```
