package job

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"guthub.com/imritik7303/boiler-plate-backend/internal/config"
)

type JobService struct{
	Client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
}

func NewJobService(logger *zerolog.Logger , cfg *config.Config) *JobService {
	redisAddr := cfg.Redis.Address

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string] int{
				"critical": 6, //higer priority queue for important eamil
				"default" : 3, //default priority for most email
				"low" :     1, //low priority for non-urgent email
			},
		},
	)

	return &JobService{
		Client: client,
		server: server,
		logger: logger,
	}
}

func (j *JobService) Start() error {
	//register task handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome , j.handleWelcomeEmailTask)

	j.logger.Info().Msg("starting background job server")
	if err := j.server.Start(mux) ; err != nil {
		return err
	}

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("stopping the background job server")
	j.server.Shutdown()
	j.Client.Close()
}