FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build employee_exercise/src/cmd/server

CMD [ "./server" ]