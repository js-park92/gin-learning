package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"gin-learning/internal/config"
	"gin-learning/internal/logger"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	logger.Init(cfg.GinMode)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialise migrator")
	}
	defer m.Close()

	if len(os.Args) < 2 {
		log.Fatal().Msg("usage: migrate <up|down> [steps]")
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal().Err(err).Msg("migrate up failed")
		}
		log.Info().Msg("migrations applied")

	case "down":
		steps := 1
		if len(os.Args) >= 3 {
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil || steps < 1 {
				log.Fatal().Msg("steps must be a positive integer")
			}
		}
		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal().Err(err).Msg("migrate down failed")
		}
		log.Info().Int("steps", steps).Msg("rolled back migration(s)")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal().Err(err).Msg("migrate version failed")
		}
		log.Info().Uint("version", version).Bool("dirty", dirty).Msg("current version")

	default:
		log.Fatal().Str("command", os.Args[1]).Msg("unknown command — use up, down, or version")
	}
}
