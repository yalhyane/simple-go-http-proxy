# Simple HTTP Proxy

This is a simple HTTP proxy project built in Golang. This was created for learning purposes to understand how HTTP proxies work.

## How to Build and Run

### Prerequisites

Before building and running the HTTP proxy, make sure you have the following installed on your system:

- Go programming language (version 1.19 or later)

### Clone the Repository
```bash
git clone https://github.com/yalhyane/simple-http-proxy
cd simple-http-proxy
go build
```

### Run the program
```shell
./simple-http-proxy --addr="127.0.0.1:8889"
```
### Run using docker
```shell
cd simple-http-proxy
docker build -t simple-go-http-proxy .
docker run --rm -it -p 8889:8889 simple-go-http-proxy
```

### Test the proxy
```shell
curl -x http://127.0.0.1:8889 http://api.ipify.org
```
