package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/StillN0THIM/sorter-up/internal/inference"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type PredictRequest struct {
	ModelName string `form:"model_name" binding:"required"`
	Version   string `form:"version"    binding:"required"`
}

type PredictResponse struct {
	RequestID string      `json:"request_id"`
	Model     string      `json:"model"`
	Version   string      `json:"version"`
	Result    interface{} `json:"result"`
	LatencyMs float64     `json:"latency_ms"`
}

func Predict(pool *pgxpool.Pool, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PredictRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//parse uploded file
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image file required"})
			return
		}
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		}
		defer f.Close()

		start := time.Now()

		result, err := inference.Run(c.Request.Context(), rdb, req.ModelName, req.Version, f)
		if err != nil {
			InferenceRequests.WithLabelValues(req.ModelName, req.Version, "errpr").Inc()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		latency := float64(time.Since(start).Milliseconds())

		//emit metrices
		InferenceRequests.WithLabelValues(req.ModelName, req.Version, "success").Inc()
		InferenceLatency.WithLabelValues(req.ModelName, req.Version).Observe(latency)

		//log in db
		go logInference(pool, req.ModelName, req.Version, latency, result)

		c.JSON(http.StatusOK, PredictResponse{
			RequestID: uuid.NewString(),
			Model:     req.ModelName,
			Version:   req.Version,
			Result:    result,
			LatencyMs: latency,
		})
	}
}

// timeout function
func logInference(pool *pgxpool.Pool, model, version string, latencyMs float64, result interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `INSERT INTO inference_log (id ,model_name, version,latency_ms,status, created_at) VALUES ($1,$2,$3,$4,'success',NOW())`, uuid.NewString(), model, version, latencyMs)

	if err != nil {
		_ = err
	}
}
