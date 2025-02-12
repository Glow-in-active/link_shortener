FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

WORKDIR /app/cmd

RUN go build -o main .

EXPOSE 8080

CMD sh -c "until nc -z -v -w30 postgres_db 5432; do echo 'Waiting for PostgreSQL...'; sleep 1; done && \
           until nc -z -v -w30 redis_db 6379; do echo 'Waiting for Redis...'; sleep 1; done && \
           ./main"