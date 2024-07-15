# go-microservices-boilerplate

A great starting point for building backend golang service with RESTful APIs/RPC in Go using negroni, and sqlx for connecting to a PostgreSQL database. The implementation follows on my confortamble working way and idioms way of writing Golang.

## Go Version
```
golang 1.12 >
```
## Run
```
make gbt
```
## Run BG-Process
```
make bg
```

## Folder Structure
```
├── cmd
│   ├── gbt
│   │   ├── gbt
│   │   ├── gbt_bg
│   │   └── main.go
│   ├── gbt_bg
│   │   └── main.go
│   └── init.go
├── configs
│   └── etc
│       └── gbt
│           └── gbt.development.ini
├── internal
│   ├── basic_module
│   │   ├── controller
│   │   │   └── controller.go
│   │   └── core
│   │       └── core.go
│   └── gbt_employee
│       ├── controller
│       │   └── controller.go
│       └── core
│           └── core.go
├── pkg
│   ├── consumer
│   │   ├── ping
│   │   │   └── ping_handler.go
│   │   └── consumer.go
│   └── http
│       ├── handler.go
│       └── routes.go
├── util
│   ├── cache
│   │   └── client
│   │       └── redis.go
│   ├── config
│   │   └── config.go
│   ├── cors
│   │   └── cors.go
│   ├── database
│   │   └── client
│   │       └── postgres.go
│   ├── middleware
│   │   └── middleware.go
│   ├── pool_worker
│   │   └── pool_worker.go
│   ├── response
│   │   └── response.go
│   ├── serve
│   │   └── serve.go
│   └── state
│       └── state.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Contributing
This project is open for contributions and suggestions. If you have an idea for a new feature or a bug fix, don't hesitate to open a pull request