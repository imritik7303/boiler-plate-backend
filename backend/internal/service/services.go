package service

import (
	"guthub.com/imritik7303/boiler-plate-backend/internal/lib/job"
	"guthub.com/imritik7303/boiler-plate-backend/internal/repository"
	"guthub.com/imritik7303/boiler-plate-backend/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Respositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}