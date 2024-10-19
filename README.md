# Tainha API Gateway

_need a logo_

**Tainha** is an open-source API Gateway developed in Go (Golang), inspired by the rich culture of Florianópolis. Designed to be an affordable and efficient solution, Tainha simplifies routing HTTP requests to various backend services with ease and flexibility.

## 🌟 Features

- [x] **Flexible Routing:** Direct requests to multiple backend services based on defined routes.
- [x] **Simple Configuration:** Use YAML files to configure routes and mappings intuitively.
- [ ] **Efficient Reverse Proxy:** Implemented with `httputil.ReverseProxy` for optimized performance.
- [x] **Dynamic Parameters Support:** Handle routes with dynamic parameters like `{companyId}`.
- [ ] **Customizable Middleware:** Easily add functionalities such as authentication, logging, and rate limiting.
- [ ] **Scalability:** Designed to handle a large volume of requests efficiently.
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
