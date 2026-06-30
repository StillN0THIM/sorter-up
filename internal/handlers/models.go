package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Model struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	TaskType    string    `json:"task_type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type RegisterModelRequest struct {
	Name        string `json:"name"        binding:"required"`
	TaskType    string `json:"task_type"   binding:"required"`
	Description string `json:"description"`
}

func ListModels(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancle()

		rows, err := pool.Query(ctx, `
			SELECT id, name, task_type, description, created_at
            FROM models
            ORDER BY created_at DESC
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch models"})
			return
		}
		defer rows.Close()

		var models []Model
		for rows.Next() {
			var m Model
			if err := rows.Scan(&m.ID, &m.Name, &m.TaskType, &m.Description, &m.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed"})
				return
			}
			models = append(models, m)
		}
		c.JSON(http.StatusOK, gin.H{"models": models})
	}
}

func GetModel(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancle := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancle()

		var m Model
		err := pool.QueryRow(ctx, `
			SELECT id, name, task_type, description, created_at
            FROM models WHERE id = $1
		`, id).Scan(&m.ID, &m.Name, &m.TaskType, &m.Description, &m.CreatedAt)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
		}
		c.JSON(http.StatusOK, m)
	}
}

func RegisterModel(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterModelRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		id := uuid.NewString()
		_, err := pool.Exec(ctx, `
            INSERT INTO models (id, name, task_type, description)
            VALUES ($1, $2, $3, $4)
        `, id, req.Name, req.TaskType, req.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register model"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}
