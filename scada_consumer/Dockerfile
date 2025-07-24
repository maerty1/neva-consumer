FROM golang:1.22.5 as builder

WORKDIR /code
COPY . .

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /code/output cmd/main.go && \
    chmod +x /code/output

FROM gcr.io/distroless/base-debian11

WORKDIR /code
COPY --from=builder /code/output /code/

USER nonroot:nonroot

CMD ["/code/output"]