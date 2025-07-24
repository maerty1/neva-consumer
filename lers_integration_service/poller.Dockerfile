FROM golang:1.22.5

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o output cmd/polling/main.go

CMD ["./output"]
