package model

type Customer struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	CustomerPhone string `json:"customer_phone"`
}
