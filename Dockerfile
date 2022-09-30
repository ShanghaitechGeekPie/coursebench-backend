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
RUN go build -ldflags="-s -w" -o backend ./cmd/coursebench-backend
RUN go build -ldflags="-s -w" -o cmd_tools ./cmd/cmd_tools

FROM alpine:latest

# Copy binary and config files from /build
# to root folder of scratch container.
COPY --from=builder ["/build/backend", "/build/cmd_tools", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/backend"]

