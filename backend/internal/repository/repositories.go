package repository

import "guthub.com/imritik7303/boiler-plate-backend/internal/server"

type Respositories struct{}

func NewRepositories(s *server.Server) *Respositories {
	return &Respositories{}
}