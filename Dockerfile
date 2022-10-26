FROM golang:1.19.2-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add curl
RUN  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /app/bin
RUN go build -o main main.go

FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate /usr/bin/migrate
COPY --from=builder /app/bin/task /usr/bin/task
COPY app.env app.env
COPY db/migration ./db/migration
COPY Taskfile.prod.yml Taskfile.yaml
EXPOSE 5000