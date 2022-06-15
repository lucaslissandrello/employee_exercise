package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"src/src/pkg/libs/employee"
	"src/src/pkg/models"
	"testing"
	"time"
)

type EmployeeManagerMock struct {
	employeeResponse *models.EmployeeResponse
	getError         error
	employeeError    employee.EmployeeError
}

func (e *EmployeeManagerMock) GetEmployees(ctx context.Context, parameters map[string]string) (*models.EmployeeResponse, error) {
	return e.employeeResponse, e.getError
}

func (e *EmployeeManagerMock) UpdateEmployeeDepartment(ctx context.Context, employeeDepartment models.EmployeeDepartment) employee.EmployeeError {
	return e.employeeError
}

func TestEmployeeController_AddEmployeeToDepartment(t *testing.T) {
	type fields struct {
		EmployeeService EmployeeManager
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		expectedResponseCode int
		expectedResponseBody *bytes.Buffer
	}{
		{
			name: "update employee department succeeds",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeError: employee.EmployeeError{},
				},
			},
			args: args{r: mockUpdateRequest(models.EmployeeDepartmentRequest{
				EmployeeNumber: 10002,
				Department:     "d006",
				FromDate:       "1996-08-04",
				ToDate:         "1996-08-09",
			})},
			expectedResponseBody: statusOkUpdatedResult(),
			expectedResponseCode: http.StatusOK,
		},
		{
			name: "update employee returns Bad request whit a wrong body type",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{},
			},
			args:                 args{r: mockUpdateRequest("wrong body")},
			expectedResponseBody: badRequestUpdatedResultWrongBody(),
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name: "update employee returns Bad request whit invalid date ranges",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{},
			},
			args: args{r: mockUpdateRequest(models.EmployeeDepartmentRequest{
				EmployeeNumber: 10002,
				Department:     "d006",
				FromDate:       "1996-08-04",
				ToDate:         "1996-08-03",
			})},
			expectedResponseBody: badRequestUpdatedResultDatesRange(),
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name: "update employee returns Bad request whit invalid from_date",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{},
			},
			args: args{r: mockUpdateRequest(models.EmployeeDepartmentRequest{
				EmployeeNumber: 10002,
				Department:     "d006",
				FromDate:       "invalid from date",
				ToDate:         "1996-08-03",
			})},
			expectedResponseBody: badRequestUpdatedResultInvalidFromDate(),
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name: "update employee returns Bad request whit invalid to_date",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{},
			},
			args: args{r: mockUpdateRequest(models.EmployeeDepartmentRequest{
				EmployeeNumber: 10002,
				Department:     "d006",
				FromDate:       "1996-08-03",
				ToDate:         "invalid to_date",
			})},
			expectedResponseBody: badRequestUpdatedResultInvalidToDate(),
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name: "update employee returns Not found when the employee is not found",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeError: employee.EmployeeError{
						Error:              sql.ErrNoRows,
						ResponseStatusCode: http.StatusNotFound,
						ErrorMessage:       "employee not found",
					},
				},
			},
			args: args{r: mockUpdateRequest(models.EmployeeDepartmentRequest{
				EmployeeNumber: 1,
				Department:     "d006",
				FromDate:       "1996-08-03",
				ToDate:         "1996-08-04",
			})},
			expectedResponseBody: employeeNotFoundUpdatedResult(),
			expectedResponseCode: http.StatusNotFound,
		},
		{
			name: "update employee returns Not found when the department is not found",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeError: employee.EmployeeError{
						Error:              sql.ErrNoRows,
						ResponseStatusCode: http.StatusNotFound,
						ErrorMessage:       "department not found",
					},
				},
			},
			args: args{r: mockUpdateRequest(models.EmployeeDepartmentRequest{
				EmployeeNumber: 10002,
				Department:     "department",
				FromDate:       "1996-08-03",
				ToDate:         "1996-08-04",
			})},
			expectedResponseBody: departmentNotFoundUpdatedResult(),
			expectedResponseCode: http.StatusNotFound,
		},
		{
			name: "update employee returns internal server error when database queries fail",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeError: employee.EmployeeError{
						Error:              sql.ErrNoRows,
						ResponseStatusCode: http.StatusInternalServerError,
						ErrorMessage:       "error in database",
					},
				},
			},
			args: args{r: mockUpdateRequest(models.EmployeeDepartmentRequest{
				EmployeeNumber: 10002,
				Department:     "d006",
				FromDate:       "1996-08-03",
				ToDate:         "1996-08-04",
			})},
			expectedResponseBody: internalServerErrorUpdateResult(),
			expectedResponseCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employeeController := &EmployeeController{
				EmployeeService: tt.fields.EmployeeService,
			}
			rr := httptest.NewRecorder()
			employeeController.AddEmployeeToDepartment(rr, tt.args.r)

			assert.Equal(t, tt.expectedResponseCode, rr.Code)
			assert.Equal(t, tt.expectedResponseBody, rr.Body)
		})
	}
}

func TestEmployeeController_GetEmployees(t *testing.T) {
	type fields struct {
		EmployeeService EmployeeManager
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		expectedResponseCode int
		expectedResponseBody *bytes.Buffer
	}{
		{
			name: "get employees with all url parameters succeeds",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeResponse: mockEmployeesResponse(),
					getError:         nil,
				},
			},
			args: args{
				request: mockRequest(),
			},
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: statusOkExpectedBody(),
		},
		{
			name: "get employees without limit and page parameters succeeds",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeResponse: mockEmployeesResponse(),
					getError:         nil,
				},
			},
			args: args{
				request: mockRequestWithoutLimitAndPage(),
			},
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: statusOkExpectedBody(),
		},
		{
			name: "get employees with wrong limit parameter returns bad request",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeResponse: nil,
					getError:         nil,
				},
			},
			args: args{
				request: mockRequestWithWrongLimitParameter(),
			},
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: badRequestWrongLimitParameterExpectedBody(),
		},
		{
			name: "get employees with wrong page parameter returns bad request",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeResponse: nil,
					getError:         nil,
				},
			},
			args: args{
				request: mockRequestWithWrongPageParameter(),
			},
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: badRequestWrongPageParameterExpectedBody(),
		},
		{
			name: "get employees fails getting data from database, returns internal server error",
			fields: fields{
				EmployeeService: &EmployeeManagerMock{
					employeeResponse: nil,
					getError:         errors.New("error getting data from database"),
				},
			},
			args: args{
				request: mockRequest(),
			},
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: internalServerErrorExpectedBody(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EmployeeController{
				EmployeeService: tt.fields.EmployeeService,
			}

			rr := httptest.NewRecorder()
			e.GetEmployees(rr, tt.args.request)

			assert.Equal(t, tt.expectedResponseCode, rr.Code)
			assert.Equal(t, tt.expectedResponseBody, rr.Body)

		})
	}
}

func mockEmployeesResponse() *models.EmployeeResponse {
	employees := []models.Employee{mockEmployee()}
	return &models.EmployeeResponse{Employees: employees, Total: 1}
}

func mockEmployee() models.Employee {
	return models.Employee{
		EmployeeNumber: 1,
		FirstName:      "Lucas",
		LastName:       "Lissandrello",
		Gender:         "M",
		BirthDate:      time.Date(1994, 11, 8, 7, 30, 00, 0, time.UTC),
		HireDate:       time.Date(2022, 06, 20, 15, 00, 00, 0, time.UTC),
	}
}

func mockRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/employees?orderBy=emp_no&order=asc&limit=5&page=1", nil)
	return request
}

func mockRequestWithoutLimitAndPage() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/employees?orderBy=emp_no&order=asc", nil)
	return request
}

func mockRequestWithWrongLimitParameter() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/employees?orderBy=emp_no&order=asc&limit=asdf&page=1", nil)
	return request
}

func mockRequestWithWrongPageParameter() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/employees?orderBy=emp_no&order=asc&limit=5&page=asdf", nil)
	return request
}

func statusOkExpectedBody() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"total":1,"page":1,"employees":[{"emp_no":1,"birth_date":"1994-11-08T07:30:00Z","first_name":"Lucas","last_name":"Lissandrello","gender":"M","hire_date":"2022-06-20T15:00:00Z"}]}`))
}

func badRequestWrongLimitParameterExpectedBody() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"bad request, wrong limit parameter"}`))
}

func badRequestWrongPageParameterExpectedBody() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"bad request, wrong page parameter"}`))
}

func internalServerErrorExpectedBody() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"internal server error"}`))
}

func mockUpdateRequest(employeeBody interface{}) *http.Request {
	body, _ := json.Marshal(employeeBody)
	request, _ := http.NewRequest(http.MethodPost, "/employees_department", bytes.NewBuffer(body))
	return request
}

func statusOkUpdatedResult() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"employee's department updated successfully"}`))
}

func badRequestUpdatedResultWrongBody() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"bad request, wrong request body"}`))
}

func badRequestUpdatedResultDatesRange() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"bad request, wrong dates range"}`))
}

func badRequestUpdatedResultInvalidFromDate() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"bad request, wrong from_date parameter"}`))
}

func badRequestUpdatedResultInvalidToDate() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"bad request, wrong to_date parameter"}`))
}

func employeeNotFoundUpdatedResult() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"employee not found"}`))
}

func departmentNotFoundUpdatedResult() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"department not found"}`))
}

func internalServerErrorUpdateResult() *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"message":"error in database"}`))
}
