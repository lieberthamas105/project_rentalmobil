package controller

import (
	"net/http"
	"project_rentalmobil/middleware"
	"project_rentalmobil/model"
	"project_rentalmobil/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Struktur CustomerController
type CustomerController struct {
	useCase        usecase.CustomerUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

// Fungsi untuk mendefinisikan route
func (c *CustomerController) Route() {
	c.rg.POST("/customers", c.authMiddleware.RequireToken("admin"), c.createNewCustomer)
	c.rg.GET("/customers", c.authMiddleware.RequireToken("admin", "user"), c.getAllCustomers)
	c.rg.GET("/customers/:id", c.authMiddleware.RequireToken("admin", "user"), c.getCustomerById)
	c.rg.PUT("/customers", c.authMiddleware.RequireToken("admin"), c.updateCustomerById)
	c.rg.DELETE("/customers/:id", c.authMiddleware.RequireToken("admin"), c.deleteCustomerById)
}

// Handler untuk mendapatkan semua pelanggan
func (c *CustomerController) getAllCustomers(ctx *gin.Context) {
	customers, err := c.useCase.GetAllCustomers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve customers data"})
		return
	}

	if len(customers) > 0 {
		ctx.JSON(http.StatusOK, customers)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "list customer is empty"})
}

// Handler untuk mendapatkan pelanggan berdasarkan ID
func (c *CustomerController) getCustomerById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	customer, err := c.useCase.GetCustomerById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get customer by ID"})
		return
	}

	ctx.JSON(http.StatusOK, customer)
}

// Handler untuk memperbarui pelanggan berdasarkan ID
func (c *CustomerController) updateCustomerById(ctx *gin.Context) {
	var payload model.Customer

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	customer, err := c.useCase.UpdateCustomerById(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, customer)
}

// Handler untuk menghapus pelanggan berdasarkan ID
func (c *CustomerController) deleteCustomerById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.useCase.DeleteCustomerById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// Handler untuk createNewCustomer
func (c *CustomerController) createNewCustomer(ctx *gin.Context) {
	var payload model.Customer

	// Bind JSON request body ke struct model.Customer
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Panggil usecase untuk menambahkan pelanggan
	customer, err := c.useCase.CreateNewCustomer(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create customer data"})
		return
	}

	// Berikan respons berhasil dengan data pelanggan
	ctx.JSON(http.StatusCreated, customer)
}

func NewCustomerController(useCase usecase.CustomerUsecase, rg *gin.RouterGroup, am middleware.AuthMiddleware) *CustomerController {
	return &CustomerController{useCase: useCase, rg: rg, authMiddleware: am}
}
