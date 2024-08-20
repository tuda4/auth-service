package main

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuda4/mb-backend/api"
	config "github.com/tuda4/mb-backend/configs"
	db "github.com/tuda4/mb-backend/db/sqlc"
	"github.com/tuda4/mb-backend/mail"
	"github.com/tuda4/mb-backend/worker"

	"github.com/rs/zerolog/log"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Error().Err(err)
	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to store")
	}

	err = server.Start(config.ADDRESS)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}

	go runTaskProcessor(config, redisOpt, store)

	log.Fatal().Msg("server start at::: " + config.ADDRESS)
}

func runTaskProcessor(config config.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Fatal().Msg("run task processor successfully")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Msg("cannot start redis processor")
	}
}
