package models

import "time"

type Employee struct {
	EmployeeNumber int       `json:"emp_no"`
	BirthDate      time.Time `json:"birth_date"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Gender         string    `json:"gender"`
	HireDate       time.Time `json:"hire_date"`
	Department     string    `json:"department"`
}

type EmployeeResponse struct {
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	Employees []Employee `json:"employees"`
}
