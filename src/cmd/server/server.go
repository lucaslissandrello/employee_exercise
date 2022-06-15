package main

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"src/src/pkg/controllers"
	"src/src/pkg/libs/database"
	"src/src/pkg/libs/employee"
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
