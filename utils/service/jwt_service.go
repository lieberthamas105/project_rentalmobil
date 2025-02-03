package service

import (
	"fmt"
	"project_rentalmobil/config"
	"project_rentalmobil/model"
	modelutil "project_rentalmobil/utils/model_util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(user model.UserCredential) (string, error)
	VerifyToken(tokenString string) (modelutil.JwtPayloadClaim, error)
}

type jwtService struct {
	cfg config.TokenConfig
}

func (j *jwtService) CreateToken(user model.UserCredential) (string, error) {
	tokenKey := j.cfg.JwtSignatureKey

	claims := modelutil.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessTokenLifeTime)),
		},
		UserId: strconv.Itoa(int(user.Id)),
		Role:   user.Role,
	}

	jwtNewClaim := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwtService) VerifyToken(tokenString string) (modelutil.JwtPayloadClaim, error) {
	tokenParse, err := jwt.ParseWithClaims(tokenString, &modelutil.JwtPayloadClaim{}, func(t *jwt.Token) (interface{}, error) {
		return j.cfg.JwtSignatureKey, nil
	})

	if err != nil {
		return modelutil.JwtPayloadClaim{}, err
	}

	claim, ok := tokenParse.Claims.(*modelutil.JwtPayloadClaim)
	if !ok {
		return modelutil.JwtPayloadClaim{}, fmt.Errorf("error claim")
	}

	return *claim, nil
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{cfg: cfg}
}
