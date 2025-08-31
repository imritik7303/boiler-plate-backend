package service

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"guthub.com/imritik7303/boiler-plate-backend/internal/server"
)

type AuthService struct {
	server *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	clerk.SetKey(s.Config.Auth.SecretKey)
	return &AuthService{
		server: s,
	}
}