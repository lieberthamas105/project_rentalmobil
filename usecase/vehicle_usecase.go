package usecase

import (
	"fmt"
	"project_rentalmobil/model"
	"project_rentalmobil/repository"
)

// Struktur vehicleUsecase
type vehicleUsecase struct {
	repo repository.VehicleRepository
}

// Interface VehicleUsecase
type VehicleUsecase interface {
	CreateNewVehicle(vehicle model.Vehicle) (model.Vehicle, error)
	GetAllVehicles() ([]model.Vehicle, error)
	GetVehicleById(id int) (model.Vehicle, error)
	UpdateVehicleById(vehicle model.Vehicle) (model.Vehicle, error)
	DeleteVehicleById(id int) error
}

// Implementasi Metode CreateNewVehicle
func (v *vehicleUsecase) CreateNewVehicle(vehicle model.Vehicle) (model.Vehicle, error) {
	return v.repo.CreateNewVehicle(vehicle)
}

// Implementasi Metode GetAllVehicles
func (v *vehicleUsecase) GetAllVehicles() ([]model.Vehicle, error) {
	return v.repo.GetAllVehicles()
}

// Implementasi Metode GetVehicleById
func (v *vehicleUsecase) GetVehicleById(id int) (model.Vehicle, error) {
	return v.repo.GetVehicleById(id)
}

// Implementasi Metode UpdateVehicleById
func (v *vehicleUsecase) UpdateVehicleById(vehicle model.Vehicle) (model.Vehicle, error) {
	_, err := v.repo.GetVehicleById(vehicle.ID)
	if err != nil {
		return model.Vehicle{}, fmt.Errorf("vehicle with ID %d not found", vehicle.ID)
	}
	return v.repo.UpdateVehicleById(vehicle)
}

// Implementasi Metode DeleteVehicleById
func (v *vehicleUsecase) DeleteVehicleById(id int) error {
	_, err := v.repo.GetVehicleById(id)
	if err != nil {
		return fmt.Errorf("vehicle with ID %d not found", id)
	}
	return v.repo.DeleteVehicleById(id)
}

// Constructor untuk VehicleUsecase
func NewVehicleUsecase(repo repository.VehicleRepository) VehicleUsecase {
	return &vehicleUsecase{repo: repo}
}
