FROM golang:alpine

ENV APIKEY="1fddfa64-1920-4c60-afa9-2675de6617bd"\
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /app

RUN cp /build/main .

EXPOSE 8080

CMD ["/app/main"]