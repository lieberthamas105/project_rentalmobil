package main

import (
	"fmt"
	"project_rentalmobil/config"
	"project_rentalmobil/controller"
	"project_rentalmobil/middleware"
	"project_rentalmobil/repository"
	"project_rentalmobil/usecase"
	"project_rentalmobil/utils/service"

	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type Server struct {
	vehicleUC  usecase.VehicleUsecase
	userUC     usecase.UserUseCase
	authUC     usecase.AuthenticationUseCase
	customerUC usecase.CustomerUsecase
	employeeUC usecase.EmployeeUsecase
	jwtService service.JwtService
	engine     *gin.Engine
	host       string
	db         *sql.DB
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)

	controller.NewVehicleController(s.vehicleUC, rg, authMiddleware).Route()
	controller.NewUserController(s.userUC, rg, authMiddleware).Route()
	controller.NewAuthController(s.authUC, rg).Route()
	controller.NewCustomerController(s.customerUC, rg, authMiddleware).Route()
	controller.NewEmployeeController(s.employeeUC, rg, authMiddleware).Route()
}

func (s *Server) Run() {
	s.initRoute()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := s.engine.Run(s.host); err != nil {
			panic(fmt.Errorf("server not running on host %s, because error %v", s.host, err.Error()))
		}
	}()

	<-sigs
	fmt.Println("Shutting down server...")
	if err := s.db.Close(); err != nil {
		fmt.Println("Error closing database connection:", err)
	} else {
		fmt.Println("Database connection closed successfully")
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.Username, cfg.DBConfig.Password, cfg.DBConfig.Database)

	db, err := sql.Open(cfg.DBConfig.Driver, dsn)
	if err != nil {
		panic("connection error: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	vehicleRepo := repository.NewVehicleRepository(db)
	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)

	vehicleUseCase := usecase.NewVehicleUsecase(vehicleRepo)
	userUseCase := usecase.NewUserUsecase(userRepo)
	customerUseCase := usecase.NewCustomerUsecase(customerRepo)
	employeeUseCase := usecase.NewEmployeeUsecase(employeeRepo)

	jwtService := service.NewJwtService(cfg.TokenConfig)
	authUseCase := usecase.NewAuthenticationUseCase(userUseCase, jwtService)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.APIConfig.ApiPort)

	return &Server{
		vehicleUC:  vehicleUseCase,
		userUC:     userUseCase,
		authUC:     authUseCase,
		customerUC: customerUseCase,
		employeeUC: employeeUseCase,
		jwtService: jwtService,
		engine:     engine,
		host:       host,
		db:         db,
	}
}
