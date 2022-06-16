package employee

import (
	"context"
	"database/sql"
	"employee_exercise/pkg/models"
	"fmt"
	"github.com/google/logger"
	"net/http"
)

type EmployeeService struct {
	EmployeeManager *sql.DB
}

type EmployeeError struct {
	Error              error
	ResponseStatusCode int
	ErrorMessage       string
}

func (e *EmployeeService) GetEmployees(ctx context.Context, parameters map[string]string) (*models.EmployeeResponse, error) {
	var employees []models.Employee
	query := fmt.Sprintf("SELECT e.emp_no, e.birth_date, e.first_name, e.last_name, e.gender, e.hire_date, d.dept_name "+
		"FROM employees e JOIN dept_emp de ON e.emp_no= de.emp_no JOIN departments d on de.dept_no = d.dept_no"+
		" ORDER BY %s %s LIMIT %s OFFSET %s", parameters["order_by_column"], parameters["order"], parameters["limit"], parameters["offset"])
	stmt, err := e.EmployeeManager.PrepareContext(ctx, query)
	if err != nil {
		logger.Errorf("error preparing sql select query: %v", err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Errorf("error executing sql select query: %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		employee := models.Employee{}
		err = rows.Scan(
			&employee.EmployeeNumber,
			&employee.BirthDate,
			&employee.FirstName,
			&employee.LastName,
			&employee.Gender,
			&employee.HireDate,
			&employee.Department,
		)
		if err != nil {
			logger.Errorf("error scanning sql select query: %v", err)
			return nil, err
		}

		employees = append(employees, employee)
	}

	total, totalError := e.getTotalEmployees(ctx)
	if totalError != nil {
		logger.Errorf("error getting total quantity from employees table: %v", err)
		return nil, totalError
	}

	employeesResponse := models.EmployeeResponse{
		Total:     total,
		Employees: employees,
	}

	return &employeesResponse, nil
}

func (e *EmployeeService) UpdateEmployeeDepartment(ctx context.Context, employeeDepartment models.EmployeeDepartment) EmployeeError {
	_, getError := e.getEmployeeByID(ctx, employeeDepartment.EmployeeNumber)
	if getError.Error != nil {
		return getError
	}

	_, departmentError := e.getDepartmentByID(ctx, employeeDepartment.Department)
	if departmentError.Error != nil {
		return departmentError
	}

	_, employeeDepartmentError := e.getEmployeeDepartment(ctx, employeeDepartment.EmployeeNumber)
	if employeeDepartmentError.Error != nil {
		if employeeDepartmentError.Error == sql.ErrNoRows {
			createError := e.createEmployeeDepartment(ctx, employeeDepartment)
			if createError.Error != nil {
				return createError
			}
			return EmployeeError{}
		}
		return employeeDepartmentError
	}

	updateError := e.updateEmployeeDepartment(ctx, employeeDepartment)
	if updateError.Error != nil {
		return updateError
	}

	return EmployeeError{}
}

func (e *EmployeeService) getEmployeeByID(ctx context.Context, employeeID int) (*models.Employee, EmployeeError) {
	var employee models.Employee
	query := fmt.Sprintf("SELECT * FROM employees WHERE emp_no=%d", employeeID)
	stmt, err := e.EmployeeManager.PrepareContext(ctx, query)
	if err != nil {
		logger.Errorf("error preparing sql select query for employee: %d, %v", employeeID, err)
		return nil, EmployeeError{
			Error:              err,
			ResponseStatusCode: http.StatusInternalServerError,
			ErrorMessage:       "error preparing sql select query",
		}
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx)
	err = row.Scan(
		&employee.EmployeeNumber,
		&employee.BirthDate,
		&employee.FirstName,
		&employee.LastName,
		&employee.Gender,
		&employee.HireDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Infof("employee not found: %d", employeeID)
			return nil, EmployeeError{
				Error:              err,
				ResponseStatusCode: http.StatusNotFound,
				ErrorMessage:       "employee not found",
			}
		} else {
			logger.Errorf("error scanning sql select query for employee: %d, %v", employeeID, err)
			return nil, EmployeeError{
				Error:              err,
				ResponseStatusCode: http.StatusInternalServerError,
				ErrorMessage:       "error scanning sql select query",
			}
		}
	}

	return &employee, EmployeeError{}
}

func (e *EmployeeService) getDepartmentByID(ctx context.Context, departmentID string) (*models.Department, EmployeeError) {
	var department models.Department
	query := fmt.Sprintf("SELECT * FROM departments WHERE dept_no='%s'", departmentID)
	stmt, err := e.EmployeeManager.PrepareContext(ctx, query)
	if err != nil {
		logger.Errorf("error preparing sql select query for department: %s, %v", departmentID, err)
		return nil, EmployeeError{
			Error:              err,
			ResponseStatusCode: http.StatusInternalServerError,
			ErrorMessage:       "error preparing sql select query",
		}
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx)
	err = row.Scan(
		&department.DepartmentNumber,
		&department.DepartmentName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Infof("department not found: %s", departmentID)
			return nil, EmployeeError{
				Error:              err,
				ResponseStatusCode: http.StatusNotFound,
				ErrorMessage:       "department not found",
			}
		} else {
			logger.Errorf("error scanning sql select query for department: %s, %v", departmentID, err)
			return nil, EmployeeError{
				Error:              err,
				ResponseStatusCode: http.StatusInternalServerError,
				ErrorMessage:       "error scanning sql select query",
			}
		}
	}

	return &department, EmployeeError{}
}

func (e *EmployeeService) getEmployeeDepartment(ctx context.Context, employeeID int) (*models.EmployeeDepartment, EmployeeError) {
	var employeeDepartment models.EmployeeDepartment
	query := fmt.Sprintf("SELECT * FROM dept_emp WHERE emp_no=%d", employeeID)
	stmt, err := e.EmployeeManager.PrepareContext(ctx, query)
	if err != nil {
		logger.Errorf("error preparing sql select query employee department for employee: %d, %v", employeeID, err)
		return nil, EmployeeError{
			Error:              err,
			ResponseStatusCode: http.StatusInternalServerError,
			ErrorMessage:       "error preparing sql select query for employee department",
		}
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx)
	err = row.Scan(
		&employeeDepartment.EmployeeNumber,
		&employeeDepartment.Department,
		&employeeDepartment.FromDate,
		&employeeDepartment.ToDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Infof("employee department not found: %d", employeeID)
			return nil, EmployeeError{
				Error:        err,
				ErrorMessage: "employee department not found",
			}
		} else {
			logger.Errorf("error scanning sql select query employee department for employee: %d, %v", employeeID, err)
			return nil, EmployeeError{
				Error:              err,
				ResponseStatusCode: http.StatusInternalServerError,
				ErrorMessage:       "error scanning sql select query for employee department",
			}
		}
	}

	return &employeeDepartment, EmployeeError{}
}

func (e *EmployeeService) updateEmployeeDepartment(ctx context.Context, employeeDepartment models.EmployeeDepartment) EmployeeError {
	query := fmt.Sprintf("UPDATE dept_emp SET dept_no='%s', from_date='%s', to_date='%s' WHERE emp_no=%d",
		employeeDepartment.Department, employeeDepartment.FromDate.Format("2006-01-02 15:04:05"),
		employeeDepartment.ToDate.Format("2006-01-02 15:04:05"), employeeDepartment.EmployeeNumber)
	stmt, err := e.EmployeeManager.PrepareContext(ctx, query)
	if err != nil {
		logger.Errorf("error preparing sql update query for employee department: %d, %v", employeeDepartment.EmployeeNumber, err)
		return EmployeeError{
			Error:              err,
			ResponseStatusCode: http.StatusInternalServerError,
			ErrorMessage:       "error preparing sql update query for employee department",
		}
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx)
	if err != nil {
		logger.Errorf("error executing sql update query for employee department: %d, %v", employeeDepartment.EmployeeNumber, err)
		return EmployeeError{
			Error:              err,
			ResponseStatusCode: http.StatusInternalServerError,
			ErrorMessage:       "error executing sql update query for employee department",
		}
	}
	rowsAffected, _ := result.RowsAffected()
	logger.Infof("employee department updated, rows affected: %d", rowsAffected)
	return EmployeeError{}
}

func (e *EmployeeService) createEmployeeDepartment(ctx context.Context, employeeDepartment models.EmployeeDepartment) EmployeeError {
	query := fmt.Sprintf("INSERT INTO dept_emp (emp_no, dept_no, from_date, to_date) VALUES (%d, '%s', '%s', '%s')",
		employeeDepartment.EmployeeNumber, employeeDepartment.Department, employeeDepartment.FromDate.Format("2006-01-02 15:04:05"), employeeDepartment.ToDate.Format("2006-01-02 15:04:05"))
	stmt, err := e.EmployeeManager.PrepareContext(ctx, query)
	if err != nil {
		logger.Errorf("error preparing sql insert query for employee department: %d, %v", employeeDepartment.EmployeeNumber, err)
		return EmployeeError{
			Error:              err,
			ResponseStatusCode: http.StatusInternalServerError,
			ErrorMessage:       "error preparing sql insert query for employee department",
		}
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx)
	if err != nil {
		logger.Errorf("error executing sql insert query for employee department: %d, %v", employeeDepartment.EmployeeNumber, err)
		return EmployeeError{
			Error:              err,
			ResponseStatusCode: http.StatusInternalServerError,
			ErrorMessage:       "error executing sql insert query for employee department",
		}
	}
	rowsAffected, _ := result.RowsAffected()
	logger.Infof("employee department created, rows affected: %d", rowsAffected)
	return EmployeeError{}
}

func (e *EmployeeService) getTotalEmployees(ctx context.Context) (int, error) {
	total := 0
	query := "SELECT COUNT(*) FROM employees"
	stmt, err := e.EmployeeManager.PrepareContext(ctx, query)
	if err != nil {
		logger.Errorf("error preparing sql select count query: %v", err)
		return 0, err
	}

	row := stmt.QueryRowContext(ctx)
	err = row.Scan(&total)
	if err != nil {
		logger.Errorf("error scanning sql count query: %v", err)
		return 0, err
	}

	return total, nil
}
