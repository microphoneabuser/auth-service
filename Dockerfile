FROM golang:1.17.2

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o auth-service ./cmd/main.go

CMD ["./auth-service"]