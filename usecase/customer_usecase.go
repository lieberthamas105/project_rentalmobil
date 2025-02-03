package usecase

import (
	"fmt"
	"project_rentalmobil/model"
	"project_rentalmobil/repository"
)

// Struktur customerUsecase
type customerUsecase struct {
	repo repository.CustomerRepository
}

// Interface CustomerUsecase
type CustomerUsecase interface {
	CreateNewCustomer(customer model.Customer) (model.Customer, error)
	GetAllCustomers() ([]model.Customer, error)
	GetCustomerById(id int) (model.Customer, error)
	UpdateCustomerById(customer model.Customer) (model.Customer, error)
	DeleteCustomerById(id int) error
}

// Implementasi Metode CreateNewCustomer
func (c *customerUsecase) CreateNewCustomer(customer model.Customer) (model.Customer, error) {
	return c.repo.CreateNewCustomer(customer)
}

// Implementasi Metode GetAllCustomers
func (c *customerUsecase) GetAllCustomers() ([]model.Customer, error) {
	return c.repo.GetAllCustomers()
}

// Implementasi Metode GetCustomerById
func (c *customerUsecase) GetCustomerById(id int) (model.Customer, error) {
	return c.repo.GetCustomerById(id)
}

// Implementasi Metode UpdateCustomerById
func (c *customerUsecase) UpdateCustomerById(customer model.Customer) (model.Customer, error) {
	_, err := c.repo.GetCustomerById(customer.ID)
	if err != nil {
		return model.Customer{}, fmt.Errorf("customer with ID %d not found", customer.ID)
	}
	return c.repo.UpdateCustomerById(customer)
}

// Implementasi Metode DeleteCustomerById
func (c *customerUsecase) DeleteCustomerById(id int) error {
	_, err := c.repo.GetCustomerById(id)
	if err != nil {
		return fmt.Errorf("customer with ID %d not found", id)
	}
	return c.repo.DeleteCustomerById(id)
}

// Constructor untuk CustomerUsecase
func NewCustomerUsecase(repo repository.CustomerRepository) CustomerUsecase {
	return &customerUsecase{repo: repo}
}
