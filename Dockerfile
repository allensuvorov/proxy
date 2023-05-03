# syntax=docker/dockerfile:1

FROM golang:latest

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-proxy

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["/docker-proxy"]






# FROM golang:1.20.3 AS builder

# WORKDIR /build

# COPY main.go .
# COPY go.mod .

# ENV CGO_ENABLED=0

# RUN go build -o proxy main.go

# FROM scratch

# COPY --from=builder /build/proxy /

# ENTRYPOINT [ "proxy" ]