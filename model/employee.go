package model

type Employee struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	EmployeePhone string `json:"employee_phone"`
}
