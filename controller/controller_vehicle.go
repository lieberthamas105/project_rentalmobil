package controller

import (
	"net/http"
	"project_rentalmobil/middleware"
	"project_rentalmobil/model"
	"project_rentalmobil/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Struktur VehicleController
type VehicleController struct {
	useCase        usecase.VehicleUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

// Fungsi untuk mendefinisikan route
func (v *VehicleController) Route() {
	v.rg.POST("/vehicles", v.authMiddleware.RequireToken("admin"), v.createNewVehicle)
	v.rg.GET("/vehicles", v.authMiddleware.RequireToken("admin", "user"), v.getAllVehicles)
	v.rg.GET("/vehicles/:id", v.authMiddleware.RequireToken("admin", "user"), v.getVehicleById)
	v.rg.PUT("/vehicles", v.authMiddleware.RequireToken("admin"), v.updateVehicleById)
	v.rg.DELETE("/vehicles/:id", v.authMiddleware.RequireToken("admin"), v.deleteVehicleById)
}

// Handler untuk mendapatkan semua kendaraan
func (v *VehicleController) getAllVehicles(c *gin.Context) {
	vehicles, err := v.useCase.GetAllVehicles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve vehicles data"})
		return
	}

	if len(vehicles) > 0 {
		c.JSON(http.StatusOK, vehicles)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "list vehicle is empty"})
}

// Handler untuk mendapatkan kendaraan berdasarkan ID
func (v *VehicleController) getVehicleById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	vehicle, err := v.useCase.GetVehicleById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get vehicle by ID"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

// Handler untuk memperbarui kendaraan berdasarkan ID
func (v *VehicleController) updateVehicleById(c *gin.Context) {
	var payload model.Vehicle

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	vehicle, err := v.useCase.UpdateVehicleById(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

// Handler untuk menghapus kendaraan berdasarkan ID
func (v *VehicleController) deleteVehicleById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := v.useCase.DeleteVehicleById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Handler untuk createNewVehicle
func (v *VehicleController) createNewVehicle(c *gin.Context) {
	var payload model.Vehicle

	// Bind JSON request body ke struct model.Vehicle
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Panggil usecase untuk menambahkan kendaraan
	vehicle, err := v.useCase.CreateNewVehicle(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create vehicle data"})
		return
	}

	// Berikan respons berhasil dengan data kendaraan
	c.JSON(http.StatusCreated, vehicle)
}

func NewVehicleController(useCase usecase.VehicleUsecase, rg *gin.RouterGroup, am middleware.AuthMiddleware) *VehicleController {
	return &VehicleController{useCase: useCase, rg: rg, authMiddleware: am}
}
