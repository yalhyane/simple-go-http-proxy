FROM golang:1.19-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY *.go ./

RUN go build -o /http-proxy-linux-amd64

EXPOSE 8889

CMD [ "/http-proxy-linux-amd64" ]
