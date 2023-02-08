package main

import (
	"challenge/internal/httpserver"
	"challenge/internal/repository"
	"challenge/pkg/env"
	"challenge/pkg/gormprovider"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
	httpPort := env.GetOrDefault("PORT", "3000")

	// Connect to database
	sqlProvider, err := gormprovider.NewSQLiteProvider()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize sqlite provider")
	}

	// Load database init
	databaseInitSql, err := os.ReadFile("database_init.sql")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database_init.sql")
	}
	err = sqlProvider.DB.Exec(string(databaseInitSql)).Error
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	server := httpserver.NewServer(
		repository.NewQuestionRepository(sqlProvider),
		repository.NewQuestionOptionRepository(sqlProvider),
	)
	err = server.Listen(fmt.Sprintf(":%s", httpPort))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start http server")
	}
}
