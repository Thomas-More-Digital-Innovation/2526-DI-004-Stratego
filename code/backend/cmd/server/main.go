// Package main is the entry point for the GoStrategy server
package main

import (
	"digital-innovation/gostrategy/internal/api"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/utils"
	"flag"
	"fmt"
)

// @title GoStrategy API
// @version 0.1.1
// @description This is the API server for the GoStrategy game.
// @termsOfService http://swagger.io/terms/

// @contact.name Sem Van Broekhoven
// @contact.url https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy
// @contact.email info@dotsem.be

// @license.name MIT
// @license.url https://opensource.org/license/mit/

// @host localhost:8080
// @BasePath /

func main() {
	if err := run(); err != nil {
		logging.Fatalf("%v", err)
	}
}

func run() error {
	defaultAddr := fmt.Sprintf(":%s", utils.GetEnv("PORT", "8080"))
	addr := flag.String("addr", defaultAddr, "Server address")
	flag.Parse()

	fmt.Println("=== GoStrategy Backend Server Running ===")

	if err := db.InitDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		if err := db.CloseDB(); err != nil {
			logging.Error("Error closing database", err)
		}
	}()

	return startServer(*addr)
}

func startServer(addr string) error {
	fmt.Printf("Starting GoStrategy Game Server on %s\n", addr)

	server := api.NewGameServer()
	if err := server.StartServer(addr); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
