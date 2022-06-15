package models

import "time"

type Department struct {
	DepartmentNumber string `json:"dept_no"`
	DepartmentName   string `json:"dept_name"`
}

type EmployeeDepartmentRequest struct {
	EmployeeNumber int    `json:"emp_no"`
	Department     string `json:"dept_no"`
	FromDate       string `json:"from_date"`
	ToDate         string `json:"to_date"`
}

type EmployeeDepartment struct {
	EmployeeNumber int       `json:"emp_no"`
	Department     string    `json:"dept_no"`
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
}
