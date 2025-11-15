FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GO111MODULE=on

CMD ["go", "run", "./main.go"]
