package controller

import (
	"net/http"
	"project_rentalmobil/middleware"
	"project_rentalmobil/model"
	"project_rentalmobil/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Struktur EmployeeController
type EmployeeController struct {
	useCase        usecase.EmployeeUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

// Fungsi untuk mendefinisikan route
func (e *EmployeeController) Route() {
	e.rg.POST("/employees", e.authMiddleware.RequireToken("admin"), e.createNewEmployee)
	e.rg.GET("/employees", e.authMiddleware.RequireToken("admin", "user"), e.getAllEmployees)
	e.rg.GET("/employees/:id", e.authMiddleware.RequireToken("admin", "user"), e.getEmployeeById)
	e.rg.PUT("/employees", e.authMiddleware.RequireToken("admin"), e.updateEmployeeById)
	e.rg.DELETE("/employees/:id", e.authMiddleware.RequireToken("admin"), e.deleteEmployeeById)
}

// Handler untuk mendapatkan semua karyawan
func (e *EmployeeController) getAllEmployees(ctx *gin.Context) {
	employees, err := e.useCase.GetAllEmployees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve employees data"})
		return
	}

	if len(employees) > 0 {
		ctx.JSON(http.StatusOK, employees)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "list employee is empty"})
}

// Handler untuk mendapatkan karyawan berdasarkan ID
func (e *EmployeeController) getEmployeeById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	employee, err := e.useCase.GetEmployeeById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get employee by ID"})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

// Handler untuk memperbarui karyawan berdasarkan ID
func (e *EmployeeController) updateEmployeeById(ctx *gin.Context) {
	var payload model.Employee

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	employee, err := e.useCase.UpdateEmployeeById(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

// Handler untuk menghapus karyawan berdasarkan ID
func (e *EmployeeController) deleteEmployeeById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := e.useCase.DeleteEmployeeById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// Handler untuk createNewEmployee
func (e *EmployeeController) createNewEmployee(ctx *gin.Context) {
	var payload model.Employee

	// Bind JSON request body ke struct model.Employee
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Panggil usecase untuk menambahkan karyawan
	employee, err := e.useCase.CreateNewEmployee(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create employee data"})
		return
	}

	// Berikan respons berhasil dengan data karyawan
	ctx.JSON(http.StatusCreated, employee)
}

func NewEmployeeController(useCase usecase.EmployeeUsecase, rg *gin.RouterGroup, am middleware.AuthMiddleware) *EmployeeController {
	return &EmployeeController{useCase: useCase, rg: rg, authMiddleware: am}
}
