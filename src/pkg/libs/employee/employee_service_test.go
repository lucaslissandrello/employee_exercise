package employee

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"src/src/pkg/models"
	"testing"
	"time"
)

func TestEmployeeService_GetEmployees_Succeeds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectQuery("emp_no", "asc", "1", "1"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRowsWithDepartment(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockCountQuery())).
		ExpectQuery().
		WithArgs().
		WillReturnRows(countRows(1))

	employeeService := &EmployeeService{EmployeeManager: db}

	employeesResponse, err := employeeService.GetEmployees(context.Background(), mockParameters())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.NotNil(t, employeesResponse)
	assert.Equal(t, expectedEmployeeResponse(), *employeesResponse)
}

func TestEmployeeService_GetEmployees_Fails_doing_select_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectQuery("emp_no", "asc", "1", "1"))).
		ExpectQuery().
		WithArgs().
		WillReturnError(errors.New("error executing query in database"))

	employeeService := &EmployeeService{EmployeeManager: db}

	employeesResponse, err := employeeService.GetEmployees(context.Background(), mockParameters())
	assert.NotNil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.Nil(t, employeesResponse)
	assert.Equal(t, errors.New("error executing query in database"), err)
}

func TestEmployeeService_GetEmployees_Fails_with_wrong_data_from_db(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectQuery("emp_no", "asc", "1", "1"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRowsWithDepartmentAndWrongData(1))

	employeeService := &EmployeeService{EmployeeManager: db}

	employeesResponse, err := employeeService.GetEmployees(context.Background(), mockParameters())
	assert.NotNil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.Nil(t, employeesResponse)
}

func TestEmployeeService_GetEmployees_Fails_preparing_select_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectQuery("emp_no", "asc", "1", "1"))).
		WillReturnError(errors.New("error preparing query in database"))

	employeeService := &EmployeeService{EmployeeManager: db}

	employeesResponse, err := employeeService.GetEmployees(context.Background(), mockParameters())
	assert.NotNil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.Nil(t, employeesResponse)
	assert.Equal(t, errors.New("error preparing query in database"), err)
}

func TestEmployeeService_GetEmployees_Fails_doing_count_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectQuery("emp_no", "asc", "1", "1"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRowsWithDepartment(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockCountQuery())).
		ExpectQuery().
		WithArgs().
		WillReturnError(errors.New("error executing count query in db"))

	employeeService := &EmployeeService{EmployeeManager: db}

	employeesResponse, err := employeeService.GetEmployees(context.Background(), mockParameters())
	assert.NotNil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.Nil(t, employeesResponse)
	assert.Equal(t, errors.New("error executing count query in db"), err)
}

func TestEmployeeService_GetEmployees_Fails_preparing_count_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectQuery("emp_no", "asc", "1", "1"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRowsWithDepartment(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockCountQuery())).
		WillReturnError(errors.New("error preparing count query in db"))

	employeeService := &EmployeeService{EmployeeManager: db}

	employeesResponse, err := employeeService.GetEmployees(context.Background(), mockParameters())
	assert.NotNil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.Nil(t, employeesResponse)
	assert.Equal(t, errors.New("error preparing count query in db"), err)
}

func TestEmployeeService_UpdateEmployeeDepartment_Succeeds_Updating(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(departmentRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeDepartmentQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeDepartmentRows(1))

	mock.ExpectPrepare(regexp.QuoteMeta(mockSqlUpdateEmployeeDepartmentQuery(employeeDepartmentUpdate))).
		ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Nil(t, updateError.Error)
}

func TestEmployeeService_UpdateEmployeeDepartment_Succeeds_Creating_new_record(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(departmentRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeDepartmentQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnError(sql.ErrNoRows)

	mock.ExpectPrepare(regexp.QuoteMeta(mockSqlInsertEmployeeDepartmentQuery(employeeDepartmentUpdate))).
		ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Nil(t, updateError.Error)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_When_Employee_not_exists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnError(sql.ErrNoRows)

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, updateError, EmployeeError{
		Error:              sql.ErrNoRows,
		ResponseStatusCode: http.StatusNotFound,
		ErrorMessage:       "employee not found",
	})
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_preparing_employee_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		WillReturnError(errors.New("error preparing sql query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, updateError, EmployeeError{
		Error:              errors.New("error preparing sql query"),
		ResponseStatusCode: http.StatusInternalServerError,
		ErrorMessage:       "error preparing sql select query",
	})
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_when_employee_has_wrong_data(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRowsWithWrongData(1))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, updateError.ResponseStatusCode, http.StatusInternalServerError)
	assert.Equal(t, "error scanning sql select query", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_When_Department_not_exists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnError(sql.ErrNoRows)

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, updateError, EmployeeError{
		Error:              sql.ErrNoRows,
		ResponseStatusCode: http.StatusNotFound,
		ErrorMessage:       "department not found",
	})
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_preparing_department_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		WillReturnError(errors.New("error preparing sql query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error preparing sql select query", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_scanning_department_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnError(errors.New("error executing sql query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error scanning sql select query", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_wrong_department_data(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error scanning sql select query", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_preparing_employee_department_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(departmentRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeDepartmentQuery(10002))).
		WillReturnError(errors.New("error preparing sql query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error preparing sql select query for employee department", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_getting_employee_department_data(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(departmentRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeDepartmentQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnError(errors.New("error executing sql query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error scanning sql select query for employee department", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_Creating_new_record_preparing_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(departmentRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeDepartmentQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnError(sql.ErrNoRows)

	mock.ExpectPrepare(regexp.QuoteMeta(mockSqlInsertEmployeeDepartmentQuery(employeeDepartmentUpdate))).
		WillReturnError(errors.New("error preparing insert sql query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error preparing sql insert query for employee department", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_Creating_new_record_executing_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(departmentRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeDepartmentQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnError(sql.ErrNoRows)

	mock.ExpectPrepare(regexp.QuoteMeta(mockSqlInsertEmployeeDepartmentQuery(employeeDepartmentUpdate))).
		ExpectExec().WithArgs().
		WillReturnError(errors.New("error preparing insert sql query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error executing sql insert query for employee department", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_Updating_preparing_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(departmentRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeDepartmentQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeDepartmentRows(1))

	mock.ExpectPrepare(regexp.QuoteMeta(mockSqlUpdateEmployeeDepartmentQuery(employeeDepartmentUpdate))).
		WillReturnError(errors.New("error preparing update query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error preparing sql update query for employee department", updateError.ErrorMessage)
}

func TestEmployeeService_UpdateEmployeeDepartment_Fails_Updating_executing_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	employeeDepartmentUpdate := models.EmployeeDepartment{
		EmployeeNumber: 10002,
		Department:     "d005",
		FromDate:       time.Date(1994, 11, 9, 7, 30, 00, 0, time.UTC),
		ToDate:         time.Date(1994, 11, 10, 7, 30, 00, 0, time.UTC),
	}

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectDepartmentQuery("d005"))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(departmentRows(1))

	mock.
		ExpectPrepare(regexp.QuoteMeta(mockSqlSelectEmployeeDepartmentQuery(10002))).
		ExpectQuery().
		WithArgs().
		WillReturnRows(employeeDepartmentRows(1))

	mock.ExpectPrepare(regexp.QuoteMeta(mockSqlUpdateEmployeeDepartmentQuery(employeeDepartmentUpdate))).
		ExpectExec().WithArgs().WillReturnError(errors.New("error executing update query"))

	employeeService := &EmployeeService{EmployeeManager: db}

	updateError := employeeService.UpdateEmployeeDepartment(context.Background(), employeeDepartmentUpdate)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, updateError.Error)
	assert.Equal(t, http.StatusInternalServerError, updateError.ResponseStatusCode)
	assert.Equal(t, "error executing sql update query for employee department", updateError.ErrorMessage)
}

func employeeRows(rowCount int) *sqlmock.Rows {
	type columnVal struct {
		name  string
		value driver.Value
	}

	rowMap := []columnVal{
		{"emp_no", "1"},
		{"birth_date", time.Date(1994, 11, 8, 7, 30, 00, 0, time.UTC)},
		{"first_name", "Lucas"},
		{"last_name", "Lissandrello"},
		{"gender", "M"},
		{"hire_date", time.Date(2022, 06, 20, 15, 00, 00, 0, time.UTC)},
	}

	var columnLabels []string
	var rowValues []driver.Value

	for _, v := range rowMap {
		columnLabels = append(columnLabels, v.name)
		rowValues = append(rowValues, v.value)
	}

	result := sqlmock.NewRows(columnLabels)

	for i := 0; i < rowCount; i++ {
		result = result.AddRow(rowValues...)
	}

	return result
}

func employeeRowsWithDepartment(rowCount int) *sqlmock.Rows {
	type columnVal struct {
		name  string
		value driver.Value
	}

	rowMap := []columnVal{
		{"emp_no", "1"},
		{"birth_date", time.Date(1994, 11, 8, 7, 30, 00, 0, time.UTC)},
		{"first_name", "Lucas"},
		{"last_name", "Lissandrello"},
		{"gender", "M"},
		{"hire_date", time.Date(2022, 06, 20, 15, 00, 00, 0, time.UTC)},
		{"dept_name", "Development"},
	}

	var columnLabels []string
	var rowValues []driver.Value

	for _, v := range rowMap {
		columnLabels = append(columnLabels, v.name)
		rowValues = append(rowValues, v.value)
	}

	result := sqlmock.NewRows(columnLabels)

	for i := 0; i < rowCount; i++ {
		result = result.AddRow(rowValues...)
	}

	return result
}

func departmentRows(rowCount int) *sqlmock.Rows {
	type columnVal struct {
		name  string
		value driver.Value
	}

	rowMap := []columnVal{
		{"dept_no", "d006"},
		{"detp_name", "Customer Service"},
	}

	var columnLabels []string
	var rowValues []driver.Value

	for _, v := range rowMap {
		columnLabels = append(columnLabels, v.name)
		rowValues = append(rowValues, v.value)
	}

	result := sqlmock.NewRows(columnLabels)

	for i := 0; i < rowCount; i++ {
		result = result.AddRow(rowValues...)
	}

	return result
}

func employeeDepartmentRows(rowCount int) *sqlmock.Rows {
	type columnVal struct {
		name  string
		value driver.Value
	}

	rowMap := []columnVal{
		{"emp_no", "10002"},
		{"dept_no", "d006"},
		{"from_date", time.Date(1994, 11, 8, 7, 30, 00, 0, time.UTC)},
		{"from_date", time.Date(1994, 11, 8, 7, 30, 00, 0, time.UTC)},
	}

	var columnLabels []string
	var rowValues []driver.Value

	for _, v := range rowMap {
		columnLabels = append(columnLabels, v.name)
		rowValues = append(rowValues, v.value)
	}

	result := sqlmock.NewRows(columnLabels)

	for i := 0; i < rowCount; i++ {
		result = result.AddRow(rowValues...)
	}

	return result
}

func employeeRowsWithWrongData(rowCount int) *sqlmock.Rows {
	type columnVal struct {
		name  string
		value driver.Value
	}

	rowMap := []columnVal{
		{"emp_no", "number"},
		{"birth_date", time.Date(1994, 11, 8, 7, 30, 00, 0, time.UTC)},
		{"first_name", 1},
		{"last_name", 2},
		{"gender", 3},
		{"hire_date", time.Date(2022, 06, 20, 15, 00, 00, 0, time.UTC)},
	}

	var columnLabels []string
	var rowValues []driver.Value

	for _, v := range rowMap {
		columnLabels = append(columnLabels, v.name)
		rowValues = append(rowValues, v.value)
	}

	result := sqlmock.NewRows(columnLabels)

	for i := 0; i < rowCount; i++ {
		result = result.AddRow(rowValues...)
	}

	return result
}

func employeeRowsWithDepartmentAndWrongData(rowCount int) *sqlmock.Rows {
	type columnVal struct {
		name  string
		value driver.Value
	}

	rowMap := []columnVal{
		{"emp_no", "number"},
		{"birth_date", time.Date(1994, 11, 8, 7, 30, 00, 0, time.UTC)},
		{"first_name", 1},
		{"last_name", 2},
		{"gender", 3},
		{"hire_date", time.Date(2022, 06, 20, 15, 00, 00, 0, time.UTC)},
		{"dept_name", 12},
	}

	var columnLabels []string
	var rowValues []driver.Value

	for _, v := range rowMap {
		columnLabels = append(columnLabels, v.name)
		rowValues = append(rowValues, v.value)
	}

	result := sqlmock.NewRows(columnLabels)

	for i := 0; i < rowCount; i++ {
		result = result.AddRow(rowValues...)
	}

	return result
}

func countRows(countResult int) *sqlmock.Rows {
	type columnVal struct {
		name  string
		value driver.Value
	}

	rowMap := []columnVal{
		{"COUNT(*)", countResult},
	}

	var columnLabels []string
	var rowValues []driver.Value

	for _, v := range rowMap {
		columnLabels = append(columnLabels, v.name)
		rowValues = append(rowValues, v.value)
	}

	result := sqlmock.NewRows(columnLabels)
	result = result.AddRow(rowValues...)

	return result
}

func expectedEmployeeResponse() models.EmployeeResponse {
	employees := []models.Employee{mockEmployee()}
	return models.EmployeeResponse{Employees: employees, Total: 1}
}

func mockEmployee() models.Employee {
	return models.Employee{
		EmployeeNumber: 1,
		FirstName:      "Lucas",
		LastName:       "Lissandrello",
		Gender:         "M",
		BirthDate:      time.Date(1994, 11, 8, 7, 30, 00, 0, time.UTC),
		HireDate:       time.Date(2022, 06, 20, 15, 00, 00, 0, time.UTC),
		Department:     "Development",
	}
}

func mockSqlSelectQuery(columnName, order, limit, offset string) string {
	sqlSelectQueryExpected := "SELECT e.emp_no, e.birth_date, e.first_name, e.last_name, e.gender, e.hire_date, d.dept_name FROM employees e JOIN dept_emp de ON e.emp_no= de.emp_no JOIN departments d on de.dept_no = d.dept_no ORDER BY " + fmt.Sprintf("%s %s", columnName, order) + " LIMIT " + limit + " OFFSET " + offset
	return sqlSelectQueryExpected
}

func mockSqlSelectEmployeeQuery(employeeID int) string {
	sqlSelectQueryExpected := fmt.Sprintf("SELECT * FROM employees WHERE emp_no=%d", employeeID)
	return sqlSelectQueryExpected
}

func mockSqlSelectDepartmentQuery(department string) string {
	sqlSelectQueryExpected := fmt.Sprintf("SELECT * FROM departments WHERE dept_no='%s'", department)
	return sqlSelectQueryExpected
}

func mockSqlSelectEmployeeDepartmentQuery(employeeID int) string {
	sqlSelectQueryExpected := fmt.Sprintf("SELECT * FROM dept_emp WHERE emp_no=%d", employeeID)
	return sqlSelectQueryExpected
}

func mockSqlUpdateEmployeeDepartmentQuery(employeeDepartment models.EmployeeDepartment) string {
	sqlSelectQueryExpected := fmt.Sprintf("UPDATE dept_emp SET dept_no='%s', from_date='%s', to_date='%s' WHERE emp_no=%d",
		employeeDepartment.Department, employeeDepartment.FromDate.Format("2006-01-02 15:04:05"),
		employeeDepartment.ToDate.Format("2006-01-02 15:04:05"), employeeDepartment.EmployeeNumber)
	return sqlSelectQueryExpected
}

func mockSqlInsertEmployeeDepartmentQuery(employeeDepartment models.EmployeeDepartment) string {
	sqlSelectQueryExpected := fmt.Sprintf("INSERT INTO dept_emp (emp_no, dept_no, from_date, to_date) VALUES (%d, '%s', '%s', '%s')",
		employeeDepartment.EmployeeNumber, employeeDepartment.Department, employeeDepartment.FromDate.Format("2006-01-02 15:04:05"), employeeDepartment.ToDate.Format("2006-01-02 15:04:05"))
	return sqlSelectQueryExpected
}

func mockCountQuery() string {
	sqlSelectQueryExpected := "SELECT COUNT(*) FROM employees"
	return sqlSelectQueryExpected
}

func mockParameters() map[string]string {
	parameters := make(map[string]string)
	parameters["limit"] = "1"
	parameters["order"] = "asc"
	parameters["order_by_column"] = "emp_no"
	parameters["offset"] = "1"

	return parameters
}
