package controllers

import (
	"context"
	"employee_exercise/pkg/libs/employee"
	"employee_exercise/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/google/logger"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type EmployeeManager interface {
	GetEmployees(ctx context.Context, parameters map[string]string) (*models.EmployeeResponse, error)
	UpdateEmployeeDepartment(ctx context.Context, employeeDepartment models.EmployeeDepartment) employee.EmployeeError
}

type EmployeeController struct {
	EmployeeService EmployeeManager
}

func (e *EmployeeController) GetEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	parameters := make(map[string]string)
	var intLimit int
	var convertError error
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		intLimit, convertError = strconv.Atoi(limit)
		if convertError != nil || intLimit < 1 {
			response["message"] = "bad request, wrong limit parameter"
			writeResponse(w, http.StatusBadRequest, response)
			return
		}
	} else {
		limit = "50"
		intLimit, _ = strconv.Atoi(limit)
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	intPage, err := strconv.Atoi(page)
	if err != nil || intPage < 1 {
		response["message"] = "bad request, wrong page parameter"
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	offset := intLimit * (intPage - 1)

	parameters["offset"] = fmt.Sprintf("%d", offset)

	order := r.URL.Query().Get("order")
	order = strings.ToLower(order)
	if order != "desc" {
		order = "asc"
	}

	parameters["order"] = order

	orderByColumn := strings.ToUpper(r.URL.Query().Get("orderBy"))
	if orderByColumn == "" {
		orderByColumn = "first_name"
	}

	parameters["order_by_column"] = orderByColumn

	parameters["limit"] = limit
	employees, err := e.EmployeeService.GetEmployees(r.Context(), parameters)
	if err != nil {
		response["message"] = "internal server error"
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	employees.Page = intPage

	writeResponse(w, http.StatusOK, employees)
}

func (e *EmployeeController) AddEmployeeToDepartment(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	var employeeDepartmentRequest models.EmployeeDepartmentRequest
	unmarshalErr := json.NewDecoder(r.Body).Decode(&employeeDepartmentRequest)
	if unmarshalErr != nil {
		logger.Errorf("error unmarshalling request body: %v", unmarshalErr)
		response["message"] = "bad request, wrong request body"
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	fromDate, toDate, errorMessage := validateDates(employeeDepartmentRequest.FromDate, employeeDepartmentRequest.ToDate)
	if errorMessage != "" {
		response["message"] = errorMessage
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	employeeDepartment := models.EmployeeDepartment{
		EmployeeNumber: employeeDepartmentRequest.EmployeeNumber,
		Department:     employeeDepartmentRequest.Department,
		FromDate:       *fromDate,
		ToDate:         *toDate,
	}

	updateError := e.EmployeeService.UpdateEmployeeDepartment(r.Context(), employeeDepartment)
	if updateError.Error != nil {
		response["message"] = updateError.ErrorMessage
		writeResponse(w, updateError.ResponseStatusCode, response)
		return
	}

	response["message"] = "employee's department updated successfully"
	writeResponse(w, http.StatusOK, response)

}

func writeResponse(w http.ResponseWriter, httpStatusCode int, response interface{}) {
	w.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(response)
	w.Write(jsonResp)
}

func validateDates(fromDateRequest, toDateRequest string) (*time.Time, *time.Time, string) {
	toDate, err := time.Parse("2006-01-02", toDateRequest)
	if err != nil {
		return nil, nil, "bad request, wrong to_date parameter"
	}

	fromDate, err := time.Parse("2006-01-02", fromDateRequest)
	if err != nil {
		return nil, nil, "bad request, wrong from_date parameter"
	}

	if toDate.Before(fromDate) {
		return nil, nil, "bad request, wrong dates range"
	}

	return &fromDate, &toDate, ""
}
