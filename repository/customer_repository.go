package repository

import (
	"database/sql"
	"project_rentalmobil/model"
)

type customerRepository struct {
	db *sql.DB
}

type CustomerRepository interface {
	CreateNewCustomer(customer model.Customer) (model.Customer, error)
	GetAllCustomers() ([]model.Customer, error)
	GetCustomerById(id int) (model.Customer, error)
	UpdateCustomerById(customer model.Customer) (model.Customer, error)
	DeleteCustomerById(id int) error
}

func (c *customerRepository) CreateNewCustomer(customer model.Customer) (model.Customer, error) {
	var customerId int

	err := c.db.QueryRow(
		"INSERT INTO mst_customer (name, address, customer_phone) VALUES ($1, $2, $3) RETURNING id",
		customer.Name, customer.Address, customer.CustomerPhone,
	).Scan(&customerId)

	if err != nil {
		return model.Customer{}, err
	}

	customer.ID = customerId
	return customer, nil
}

func (c *customerRepository) GetAllCustomers() ([]model.Customer, error) {
	var customers []model.Customer

	rows, err := c.db.Query("SELECT id, name, address, customer_phone FROM mst_customer")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var customer model.Customer

		err := rows.Scan(&customer.ID, &customer.Name, &customer.Address, &customer.CustomerPhone)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (c *customerRepository) GetCustomerById(id int) (model.Customer, error) {
	var customer model.Customer

	err := c.db.QueryRow(
		"SELECT id, name, address, customer_phone FROM mst_customer WHERE id = $1", id,
	).Scan(&customer.ID, &customer.Name, &customer.Address, &customer.CustomerPhone)

	if err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}

func (c *customerRepository) UpdateCustomerById(customer model.Customer) (model.Customer, error) {
	_, err := c.db.Exec(
		"UPDATE mst_customer SET name = $2, address = $3, customer_phone = $4 WHERE id = $1",
		customer.ID, customer.Name, customer.Address, customer.CustomerPhone,
	)

	if err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}

func (c *customerRepository) DeleteCustomerById(id int) error {
	_, err := c.db.Exec("DELETE FROM mst_customer WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}
