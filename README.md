# test-proxy
HTTP server for proxying HTTP-requests to 3rd-party services.

## Features

- Timeouts. Client requests handler is configured with time limits for: read, write, idle.
- Basic caching.The caching mechanism is based on storing the request-response pairs in a map. A duplicate request gets the response stored in the map for the same request.
- Dockerfile and Dockerfile.multistage to build docker images.
- Basic end-to-end testing. 

## Docker
To build an images: 
- $ docker build --tag docker-proxy .

To run the app: 
- $ docker run --publish 8080:8080 docker-proxy

For Multistage:
- $ docker build -t docker-proxy-ms:multistage -f Dockerfile.multistage .
- $ docker run --publish 8080:8080 docker-proxy:multistage

## cURL testing

$ curl -X POST http://localhost:8080 -H 'Content-Type: application/json' -d '{"method": "GET", "url": "https://catfact.ninja/fact", "headers": {"Authorization": "Bearer abc123"}}'