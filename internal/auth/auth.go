package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

// TODO:
// Not sure how exactly JWT here is needed, read the README.md tomorrow again.

type Claims = jwt.MapClaims

type Service struct {
	jwtKey []byte
}

func NewService(jwtKey []byte) *Service {
	return &Service{jwtKey: jwtKey}
}

func (s *Service) GenerateJWT() {
}
