## EMPLOYEE-EXERCISE

This is an exercise. It uses the datacharmer database from this repo:

https://github.com/datacharmer/test_db.

### Endpoints


  #### Get all employees

    curl --location --request GET '/employees?orderBy=emp_no&order=asc&limit=50&page=1'
  
  It returns the employees. Has the following URL parameters:

    -limit(int): used to limit the response, if it is not present, the default value will be 50.
  
    -page(int): used as offset for pagination
  
    -order(string): asc or desc, default value is "asc"
  
    -orderBy(string): column to order, default value is "first_name"



#### Update employee's department

    curl --location --request POST '/employees_department' \ --header 'Content-Type: application/json' \ --data-raw '{ "emp_no": 10002, "dept_no": "d002", "from_date": "1996-08-04", "to_date": "1996-08-07" }'

  It updates the employee's department. If the employee has no department assigned, it will add the requested department.
  
  The body request:

  
    -emp_no(int): number of employee
    
    -dept_no(string): department number
    
    -from_date(string): date from
    
    -to_date(string): date to

### Required Env Vars ###

* **MYSQL_USER**

  Ex: MYSQL_USER = user

* **MYSQL_PASSWORD**

  Ex: MYSQL_PASSWORD = password

* **DB_NAME**

  Ex: DB_NAME = database_name

* **MYSQL_HOST**

  Ex: MYSQL_HOST = database_host

* **MYSQL_PORT**

  Ex: MYSQL_PORT = 3306