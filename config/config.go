package config

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type APIConfig struct {
	ApiPort string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     []byte
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type Config struct {
	DBConfig
	APIConfig
	TokenConfig TokenConfig
}

// Fungsi untuk membaca konfigurasi dari .env
func (c *Config) readConfig() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Konversi waktu kadaluwarsa token dari string ke time.Duration
	tokenLifetimeStr := os.Getenv("ACCESS_TOKEN_LIFETIME")
	tokenLifetime, err := time.ParseDuration(tokenLifetimeStr)
	if err != nil {
		tokenLifetime = 1 * time.Hour // Default jika terjadi kesalahan parsing
	}

	// Set konfigurasi dari env
	c.DBConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Driver:   os.Getenv("DB_DRIVER"),
	}
	c.APIConfig = APIConfig{
		ApiPort: os.Getenv("API_PORT"),
	}

	c.TokenConfig = TokenConfig{
		ApplicationName:     os.Getenv("APP_NAME"),
		JwtSignatureKey:     []byte(os.Getenv("JWT_SIGNATURE_KEY")),
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: tokenLifetime,
	}

	// Validasi wajib
	if c.DBConfig.Host == "" || c.DBConfig.Port == "" || c.DBConfig.Username == "" || c.DBConfig.Password == "" || c.APIConfig.ApiPort == "" {
		return fmt.Errorf("required config is missing")
	}

	if c.TokenConfig.ApplicationName == "" || len(c.TokenConfig.JwtSignatureKey) == 0 || c.TokenConfig.JwtSigningMethod == nil {
		return fmt.Errorf("token config is invalid")
	}

	return nil
}

// Fungsi untuk mendapatkan TokenConfig
func (c *Config) GetTokenConfig() TokenConfig {
	return c.TokenConfig
}

// Fungsi untuk membuat instance Config baru
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cfg.readConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
