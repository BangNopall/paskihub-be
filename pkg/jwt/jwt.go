package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/pkg/log"
)

type JwtInterface interface {
	GenerateToken(userId uuid.UUID, payload entity.User) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, string, string, *uuid.UUID, error)
}

type JwtStruct struct {
	SecretKey   string
	ExpiredTime time.Duration
}

type Claims struct {
	Id       uuid.UUID
	Email    string
	Role     string
	ParentId *uuid.UUID `json:"parent_id,omitempty"`
	jwt.RegisteredClaims
}

var Jwt = getJwt()

func getJwt() JwtInterface {
	secretKey := env.AppEnv.JwtSecretKey
	expiredTime, err := strconv.Atoi(env.AppEnv.JwtExpireTime)
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[JWT][getJwt] failed to convert string to int")
	}

	return &JwtStruct{
		SecretKey:   secretKey,
		ExpiredTime: time.Duration(expiredTime) * time.Hour,
	}
}

func (j *JwtStruct) GenerateToken(id uuid.UUID, payload entity.User) (string, error) {
	claim := &Claims{
		Id:       id,
		Email:    payload.Email,
		Role:     string(payload.Role),
		ParentId: payload.ParentId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiredTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[JWT][GenerateToken] failed to generate token")

		return tokenString, err
	}

	return tokenString, nil
}

func (j *JwtStruct) ValidateToken(tokenString string) (uuid.UUID, string, string, *uuid.UUID, error) {
	var id uuid.UUID
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[JWT][ValidateToken] failed to validate token")

		return id, "", "", nil, err
	}

	if !token.Valid {
		log.Warn(log.LogInfo{
			"error": "invalid token",
		}, "[JWT][ValidateToken] invalid token")

		return id, "", "", nil, domain.ErrInvalidToken
	}

	id = claims.Id
	email := claims.Email
	role := claims.Role

	return id, email, role, claims.ParentId, nil
}
