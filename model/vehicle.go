package model

type Vehicle struct {
	ID           int    `json:"id"`
	BrandName    string `json:"brand_name"`
	YearReleased string `json:"year_released"`
	LicensePlate string `json:"license_plate"`
	Kilometer    int    `json:"kilometer"`
}
