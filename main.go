package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"gin-learning/internal/config"
	"gin-learning/internal/database"
	"gin-learning/internal/handler"
	"gin-learning/internal/logger"
	"gin-learning/internal/repository"
	"gin-learning/internal/router"
	"gin-learning/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	logger.Init(cfg.GinMode)

	database.Connect(cfg.Database)

	userRepo := repository.NewUserRepository(database.DB)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	r := router.New(userHandler)

	log.Info().Str("port", cfg.Port).Msg("starting server")
	r.Run(":" + cfg.Port)
}
