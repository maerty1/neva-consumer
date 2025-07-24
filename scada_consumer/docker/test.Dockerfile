FROM golang:1.22.5
 
RUN curl -sSL https://get.docker.com/ | sh
WORKDIR /code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["sh", "-c", "dockerd --log-level=error > /dev/null 2>&1 & go test ./..."]
