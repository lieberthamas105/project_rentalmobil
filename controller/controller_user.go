package controller

import (
	"net/http"
	"strconv"

	"project_rentalmobil/middleware"
	"project_rentalmobil/model"
	"project_rentalmobil/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	useCase        usecase.UserUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (b *UserController) createUser(c *gin.Context) {
	var payload model.UserCredential

	// Bind JSON payload ke struct UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Panggil useCase untuk mendaftarkan user baru
	user, err := b.useCase.RegisterNewUser(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create user"})
		return
	}

	// Kembalikan respons berhasil dengan data user yang dibuat
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (b *UserController) getAllUser(c *gin.Context) {
	users, err := b.useCase.FindAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data"})
		return
	}

	if len(users) > 0 {
		c.JSON(http.StatusOK, users)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "List user empty"})
}

func (b *UserController) getUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid ID"})
		return
	}

	// Call use case
	user, err := b.useCase.FindUserById(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get user by ID"})
		return
	}

	// Check if user is not found (e.g., user.Id == 0)
	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (b *UserController) Route() {
	b.rg.POST("/users", b.authMiddleware.RequireToken("admin"), b.createUser)
	b.rg.GET("/users", b.authMiddleware.RequireToken("admin"), b.getAllUser)
	b.rg.GET("/users/:id", b.authMiddleware.RequireToken("admin"), b.getUserById)
}

func NewUserController(useCase usecase.UserUseCase, rg *gin.RouterGroup, am middleware.AuthMiddleware) *UserController {
	return &UserController{useCase: useCase, rg: rg, authMiddleware: am}
}
