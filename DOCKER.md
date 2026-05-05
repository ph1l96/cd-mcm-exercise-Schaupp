# Dockerfile Analysis

The `Dockerfile` uses a multi-stage build with two stages:

- a `builder` stage based on `golang:1.26-alpine`
- a final runtime stage based on `alpine:3.19`

## Stage 1: Builder

```dockerfile
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./cmd/api
```

Purpose of this stage:

- provides the Go compiler and build tools
- downloads the Go module dependencies
- copies the source code into the image
- compiles the API binary as `/api-server`

Why it is structured this way:

- copying `go.mod` and `go.sum` before the rest of the source allows Docker to cache dependency download separately
- `go mod download` only needs to rerun when dependencies change
- the final `go build` produces a Linux binary for the runtime image

## Stage 2: Runtime

```dockerfile
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /api-server .

EXPOSE 8080

ENTRYPOINT ["./api-server"]
```

Purpose of this stage:

- creates a much smaller image for running the service
- installs only `ca-certificates`, which many Go applications need for TLS and HTTPS
- copies only the compiled binary from the builder stage
- exposes port `8080`
- starts the application directly

This stage does not include:

- the Go compiler
- the source code
- module caches
- build tools

That is the main advantage of a multi-stage build.

## What `CGO_ENABLED=0` Does

`CGO_ENABLED=0` disables CGO, meaning the Go compiler will not link against C libraries.

Why that matters here:

- it encourages a fully static Go binary
- the binary becomes easier to run in a minimal runtime image
- it reduces dependency on native system libraries that may not exist in the final container

Why it is important in this Dockerfile:

- the runtime image is plain Alpine and only installs `ca-certificates`
- if the binary depended on external C libraries, it could fail at runtime because those libraries are not present
- with `CGO_ENABLED=0`, the binary is more portable and better suited to a stripped-down final image

In short: `CGO_ENABLED=0` helps ensure the compiled application can run in the minimal runtime container without requiring extra libc-related dependencies.

## Final Image Size vs Single-Stage Build

### Current Multi-Stage Build

The final image contains only:

- Alpine Linux
- CA certificates
- the compiled `api-server` binary

This is relatively small because it excludes the full Go toolchain and all build-time artifacts.

### Typical Single-Stage Build

A single-stage Dockerfile would usually look like this:

```dockerfile
FROM golang:1.26-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./cmd/api
EXPOSE 8080
ENTRYPOINT ["./api-server"]
```

That image would keep everything inside one image:

- the Go compiler
- package manager metadata
- module download cache
- source code
- build artifacts
- the final binary

## Size Trade-Off

Compared to a single-stage build, the multi-stage image is usually much smaller.

Expected result:

- multi-stage build: small runtime-focused image
- single-stage build: significantly larger image because it includes the whole Go build environment

Benefits of the smaller final image:

- faster image transfer and deployment
- lower storage usage
- smaller attack surface
- less unnecessary tooling in production

# TASK 3 - CRUD Test results

#### Create 3 products
````
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name":"Product 1","price":9.99}'

{"id":1,"name":"Product 1","price":9.99}
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name":"Product 2","price":19.99}'

{"id":2,"name":"Product 2","price":19.99}
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name":"Product 3","price":29.99}'

{"id":3,"name":"Product 3","price":29.99}
````

#### List
````
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl http://localhost:8080/products                                                                                   
[{"id":1,"name":"Product 1","price":9.99},{"id":2,"name":"Product 2","price":19.99},{"id":3,"name":"Product 3","price":29.99}]
````

#### Update
````
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl -X PUT http://localhost:8080/products/1 -H "Content-Type: application/json" -d '{"name":"Product 1 v2","price":67.67}'
{"id":1,"name":"Product 1 v2","price":67.67}
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl http://localhost:8080/products                                                                                        
[{"id":1,"name":"Product 1 v2","price":67.67},{"id":2,"name":"Product 2","price":19.99},{"id":3,"name":"Product 3","price":29.99}]
````

#### Delete
````
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl -X DELETE http://localhost:8080/products/3
{"result":"success"}
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl http://localhost:8080/products            
[{"id":1,"name":"Product 1 v2","price":67.67},{"id":2,"name":"Product 2","price":19.99}]
````

#### Persistence
````
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) docker compose down
WARN[0000] /home/schauppp/Documents/MC_Master/ContinuousDelivery/exercises/cd-mcm-exercise-Schaupp/docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion 
[+] down 3/3
 ✔ Container cd-mcm-exercise-schaupp-api-1 Removed                                                                                                                                                                                                                                          0.3s
 ✔ Container cd-mcm-exercise-schaupp-db-1  Removed                                                                                                                                                                                                                                          0.3s
 ✔ Network cd-mcm-exercise-schaupp_default Removed                                                                                                                                                                                                                                          0.2s
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) docker compose up -d
WARN[0000] /home/schauppp/Documents/MC_Master/ContinuousDelivery/exercises/cd-mcm-exercise-Schaupp/docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion 
[+] up 3/3
 ✔ Network cd-mcm-exercise-schaupp_default Created                                                                                                                                                                                                                                          0.2s
 ✔ Container cd-mcm-exercise-schaupp-db-1  Healthy                                                                                                                                                                                                                                          6.1s
 ✔ Container cd-mcm-exercise-schaupp-api-1 Started                                                                                                                                                                                                                                          6.3s
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) curl http://localhost:8080/products
[{"id":1,"name":"Product 1 v2","price":67.67},{"id":2,"name":"Product 2","price":19.99}]
➜  cd-mcm-exercise-Schaupp git:(exercise/02-microservice-docker) 
````