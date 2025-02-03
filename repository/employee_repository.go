package repository

import (
	"database/sql"
	"project_rentalmobil/model"
)

type employeeRepository struct {
	db *sql.DB
}

type EmployeeRepository interface {
	CreateNewEmployee(employee model.Employee) (model.Employee, error)
	GetAllEmployees() ([]model.Employee, error)
	GetEmployeeById(id int) (model.Employee, error)
	UpdateEmployeeById(employee model.Employee) (model.Employee, error)
	DeleteEmployeeById(id int) error
}

func (e *employeeRepository) CreateNewEmployee(employee model.Employee) (model.Employee, error) {
	var employeeId int

	err := e.db.QueryRow(
		"INSERT INTO mst_employee (name, address, employee_phone) VALUES ($1, $2, $3) RETURNING id",
		employee.Name, employee.Address, employee.EmployeePhone,
	).Scan(&employeeId)

	if err != nil {
		return model.Employee{}, err
	}

	employee.ID = employeeId
	return employee, nil
}

func (e *employeeRepository) GetAllEmployees() ([]model.Employee, error) {
	var employees []model.Employee

	rows, err := e.db.Query("SELECT id, name, address, employee_phone FROM mst_employee")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var employee model.Employee

		err := rows.Scan(&employee.ID, &employee.Name, &employee.Address, &employee.EmployeePhone)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

func (e *employeeRepository) GetEmployeeById(id int) (model.Employee, error) {
	var employee model.Employee

	err := e.db.QueryRow(
		"SELECT id, name, address, employee_phone FROM mst_employee WHERE id = $1", id,
	).Scan(&employee.ID, &employee.Name, &employee.Address, &employee.EmployeePhone)

	if err != nil {
		return model.Employee{}, err
	}

	return employee, nil
}

func (e *employeeRepository) UpdateEmployeeById(employee model.Employee) (model.Employee, error) {
	_, err := e.db.Exec(
		"UPDATE mst_employee SET name = $2, address = $3, employee_phone = $4 WHERE id = $1",
		employee.ID, employee.Name, employee.Address, employee.EmployeePhone,
	)

	if err != nil {
		return model.Employee{}, err
	}

	return employee, nil
}

func (e *employeeRepository) DeleteEmployeeById(id int) error {
	_, err := e.db.Exec("DELETE FROM mst_employee WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}
