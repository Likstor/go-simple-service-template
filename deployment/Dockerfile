FROM golang:latest

WORKDIR /

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux go build -o /migrator ./cmd/migrator/migrator.go

RUN GOOS=linux go build -o /service ./cmd/service/service.go

RUN sed -i 's/\r$//' /deployment/entrypoint.sh && \
    chmod +x /deployment/entrypoint.sh

ENTRYPOINT ["/deployment/entrypoint.sh"]