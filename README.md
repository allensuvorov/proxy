# test-proxy
HTTP server for proxying HTTP-requests to 3rd-party services.

## Features

- Timeouts. Client requests handler is configured with time limits for: read, write, idle.
- Caching. Basic caching with a hash map. A duplicate request gets the response for the same request stored in the map.

## Docker
- to build an images: $ docker build --tag docker-proxy .
- to run the app : $ docker run --publish 8080:8080 docker-proxy