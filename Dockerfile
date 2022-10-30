FROM golang:1.19-alpine

WORKDIR /urls_checker

COPY . .

RUN go mod tidy

CMD go run backend/main.go

EXPOSE 4000