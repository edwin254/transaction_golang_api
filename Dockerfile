FROM golang:1.23.4-alpine

WORKDIR /app

# Install git for downloading Air
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest


COPY go.mod ./
COPY go.sum ./
RUN go mod download  && go mod verify

COPY . .
COPY .air.toml . 

RUN go build -v -o /cmd/api/main .

CMD ["/cmd/api/main"]