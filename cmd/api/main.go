package main

import (
	"log"
	"os"

	"github.com/StillN0THIM/inference-platform/internal/config"
	"github.com/StillN0THIM/inference-platform/internal/db"
	"github.com/StillN0THIM/inference-platform/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found,reading from environment")
	}

	cfg := config.Load()

	pool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Postgress connection failed: %v", err)
	}
	defer pool.Close()

	rdb := db.NewRedisClient(cfg.RedisAddr)
	defer rdb.Close()

	r := gin.Default()

	routes.Register(r, pool, rdb)

	log.Printf("server started on %s", &cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server faild to start start: %v", err)
	}
}
