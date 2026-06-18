package routes

import (
	"github.com/StillN0THIM/sorter-up/internal/handlers"
	"github.com/StillN0THIM/sorter-up/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Register(r *gin.Engine, pool *pgxpool.Pool, rdb *redis.Client) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	r.Get("/health", handlers.Health)

	r.Get("metrics", handlers.Matrices())

	v1 := r.Group("/v1")
	{
		models := v1.Group("/models")
		{
			models.GET("", handlers.ListModels(pool))
			models.POST("", handlers.RegisterModel(pool))
			models.GET("/:id", handlers.GetModels(pool))
		}

		v1.POST("/predict", handlers.Predict(pool, rdb))
		v1.POST("/ocr", handlers.OCR(pool, rdb))
		v1.POST(".detect", handlers.Detect(pool, rdb))
	}
}
