package handlers

import (
	"guthub.com/imritik7303/boiler-plate-backend/internal/server"
	"guthub.com/imritik7303/boiler-plate-backend/internal/service"
)

type Handlers struct{}

func NewHandlers(s *server.Server , services *service.Services) *Handlers {
	return &Handlers{}
}