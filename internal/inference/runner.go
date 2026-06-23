package inference

import (
	"context"
	"io"

	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, rdb *redis.Client, modelName, version string, r io.Reader) (interface{}, error) {
	return map[string]string{"status": "stub"}, nil
}
