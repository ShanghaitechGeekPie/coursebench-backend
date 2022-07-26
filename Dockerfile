FROM golang:alpine AS builder

RUN apk add git

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image
# and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o backend cmd/coursebench-backend/main.go
RUN go build -ldflags="-s -w" -o import_course cmd/import_course/main.go

FROM debian:stable-slim

# Copy binary and config files from /build
# to root folder of scratch container.
COPY --from=builder ["/build/backend", "/build/import_course", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/backend"]

