package repository

import (
	"database/sql"
	"project_rentalmobil/model"
)

type vehicleRepository struct {
	db *sql.DB
}

type VehicleRepository interface {
	CreateNewVehicle(vehicle model.Vehicle) (model.Vehicle, error)
	GetAllVehicles() ([]model.Vehicle, error)
	GetVehicleById(id int) (model.Vehicle, error)
	UpdateVehicleById(vehicle model.Vehicle) (model.Vehicle, error)
	DeleteVehicleById(id int) error
}

func (v *vehicleRepository) CreateNewVehicle(vehicle model.Vehicle) (model.Vehicle, error) {
	var vehicleId int

	err := v.db.QueryRow(
		"INSERT INTO mst_vehicle (brand_name, year_released, license_plate, kilometer) VALUES ($1, $2, $3, $4) RETURNING id",
		vehicle.BrandName, vehicle.YearReleased, vehicle.LicensePlate, vehicle.Kilometer,
	).Scan(&vehicleId)

	if err != nil {
		return model.Vehicle{}, err
	}

	vehicle.ID = vehicleId
	return vehicle, nil
}

func (v *vehicleRepository) GetAllVehicles() ([]model.Vehicle, error) {
	var vehicles []model.Vehicle

	rows, err := v.db.Query("SELECT id, brand_name, year_released, license_plate, kilometer FROM mst_vehicle")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var vehicle model.Vehicle

		err := rows.Scan(&vehicle.ID, &vehicle.BrandName, &vehicle.YearReleased, &vehicle.LicensePlate, &vehicle.Kilometer)
		if err != nil {
			return nil, err
		}

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func (v *vehicleRepository) GetVehicleById(id int) (model.Vehicle, error) {
	var vehicle model.Vehicle

	err := v.db.QueryRow(
		"SELECT id, brand_name, year_released, license_plate, kilometer FROM mst_vehicle WHERE id = $1", id,
	).Scan(&vehicle.ID, &vehicle.BrandName, &vehicle.YearReleased, &vehicle.LicensePlate, &vehicle.Kilometer)

	if err != nil {
		return model.Vehicle{}, err
	}

	return vehicle, nil
}

func (v *vehicleRepository) UpdateVehicleById(vehicle model.Vehicle) (model.Vehicle, error) {
	_, err := v.db.Exec(
		"UPDATE mst_vehicle SET brand_name = $2, year_released = $3, license_plate = $4, kilometer = $5 WHERE id = $1",
		vehicle.ID, vehicle.BrandName, vehicle.YearReleased, vehicle.LicensePlate, vehicle.Kilometer,
	)

	if err != nil {
		return model.Vehicle{}, err
	}

	return vehicle, nil
}

func (v *vehicleRepository) DeleteVehicleById(id int) error {
	_, err := v.db.Exec("DELETE FROM mst_vehicle WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func NewVehicleRepository(db *sql.DB) VehicleRepository {
	return &vehicleRepository{db: db}
}
