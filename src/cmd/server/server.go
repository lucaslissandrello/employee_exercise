package main

import (
	"employee_exercise/src/pkg/controllers"
	"employee_exercise/src/pkg/libs/database"
	"employee_exercise/src/pkg/libs/employee"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {

	envErrors := readEnv()
	if len(envErrors) != 0 {
		logger.Fatal("could not process environment: ", envErrors)
		os.Exit(1)
	}

	employeeController := controllers.EmployeeController{
		EmployeeService: &employee.EmployeeService{
			EmployeeManager: database.GetDbEngine(),
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/employees", employeeController.GetEmployees).Methods("GET")
	router.HandleFunc("/employees_department", employeeController.AddEmployeeToDepartment).Methods("POST")

	err := http.ListenAndServe(":80", router)
	if err != nil {
		logger.Fatal(nil, "Error listening on port 80 ")
		return
	}
}
