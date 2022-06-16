FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY src/cmd/server/env.go ./src/cmd/server/env.go
COPY src/cmd/server/server.go ./src/cmd/server/server.go
COPY src/pkg/controllers/employee_controller.go ./src/pkg/controllers/employee_controller.go
COPY src/pkg/controllers/employee_controller_test.go ./src/pkg/controllers/employee_controller_test.go
COPY src/pkg/libs/database/database.go ./src/pkg/libs/database/database.go
COPY src/pkg/libs/employee/employee_service.go ./src/pkg/libs/employee/employee_service.go
COPY src/pkg/libs/employee/employee_service_test.go ./src/pkg/libs/employee/employee_service_test.go
COPY src/pkg/models/department.go ./src/pkg/models/department.go
COPY src/pkg/models/employee.go ./src/pkg/models/employee.go

RUN ls src/cmd/aasdf

RUN go build -o src/cmd/aasdf/aasdf

WORKDIR /dist

RUN cp /build/main .

CMD [ "/build/server" ]