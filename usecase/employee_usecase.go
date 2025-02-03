package usecase

import (
	"fmt"
	"project_rentalmobil/model"
	"project_rentalmobil/repository"
)

// Struktur employeeUsecase
type employeeUsecase struct {
	repo repository.EmployeeRepository
}

// Interface EmployeeUsecase
type EmployeeUsecase interface {
	CreateNewEmployee(employee model.Employee) (model.Employee, error)
	GetAllEmployees() ([]model.Employee, error)
	GetEmployeeById(id int) (model.Employee, error)
	UpdateEmployeeById(employee model.Employee) (model.Employee, error)
	DeleteEmployeeById(id int) error
}

// Implementasi Metode CreateNewEmployee
func (e *employeeUsecase) CreateNewEmployee(employee model.Employee) (model.Employee, error) {
	return e.repo.CreateNewEmployee(employee)
}

// Implementasi Metode GetAllEmployees
func (e *employeeUsecase) GetAllEmployees() ([]model.Employee, error) {
	return e.repo.GetAllEmployees()
}

// Implementasi Metode GetEmployeeById
func (e *employeeUsecase) GetEmployeeById(id int) (model.Employee, error) {
	return e.repo.GetEmployeeById(id)
}

// Implementasi Metode UpdateEmployeeById
func (e *employeeUsecase) UpdateEmployeeById(employee model.Employee) (model.Employee, error) {
	_, err := e.repo.GetEmployeeById(employee.ID)
	if err != nil {
		return model.Employee{}, fmt.Errorf("employee with ID %d not found", employee.ID)
	}
	return e.repo.UpdateEmployeeById(employee)
}

// Implementasi Metode DeleteEmployeeById
func (e *employeeUsecase) DeleteEmployeeById(id int) error {
	_, err := e.repo.GetEmployeeById(id)
	if err != nil {
		return fmt.Errorf("employee with ID %d not found", id)
	}
	return e.repo.DeleteEmployeeById(id)
}

// Constructor untuk EmployeeUsecase
func NewEmployeeUsecase(repo repository.EmployeeRepository) EmployeeUsecase {
	return &employeeUsecase{repo: repo}
}
