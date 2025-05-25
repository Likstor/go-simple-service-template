FROM golang:latest

WORKDIR /

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux go build -o /migrator ./cmd/migrator/migrator.go

RUN GOOS=linux go build -o /service ./cmd/service/service.go

COPY ./deployment/entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]