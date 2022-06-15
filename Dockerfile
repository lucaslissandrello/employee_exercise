FROM golang:alpine AS build
WORKDIR /go/src/employee
COPY . .
RUN go build -o /go/bin/employee cmd/server.go

FROM scratch
COPY --from=build /go/bin/employee /go/bin/employee
ENTRYPOINT ["/go/bin/employee"]