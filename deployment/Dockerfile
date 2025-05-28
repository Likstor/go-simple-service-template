FROM golang:latest AS builder

WORKDIR /builder

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /builder/service ./cmd/service/service.go

FROM alpine:latest

WORKDIR /service

COPY --from=builder /builder/service .

COPY deployment/entrypoint.sh .

RUN sed -i 's/\r$//' entrypoint.sh && \
    chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]