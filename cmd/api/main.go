package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/oktaharis/uji-teknis-godigi/internal/config"
	"github.com/oktaharis/uji-teknis-godigi/internal/database"
	"github.com/oktaharis/uji-teknis-godigi/internal/routes"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	db := database.Connect(cfg)

	r := routes.SetupRouter(cfg, db)

	log.Printf("server running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Println("failed to start:", err)
		os.Exit(1)
	}
}
