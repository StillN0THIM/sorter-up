package loader

import (
	"context"
	"fmt"
	"go/version"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
	ort "github.com/yalue/onnxruntime_go"
)

func sessionKey(modelName, version string) string {
	return fmt.Sprintf("session:%s:%s", modelName, version)
}

func GetSession(ctx context.Context, rdb *redis.Client, modelName, version string) (*ort.DynamicAdvancedSession, error) {
	key := sessionKey(modelName, version)

	if session, ok := sessionCache.get(key); ok {
		return session, nil
	}

	session, err := loadFromDisk(modelName, version)
	if err != nil {
		return nil, err
	}

	sessionCache.set(key, session)
	rdb.Set(ctx, key, "loaded", 0)

	return session, nil
}

func loadFromDisk(modelName, version string) (*ort.DynamicAdvancedSession, error) {
	modelsDir := os.Getenv("MODELS_DIR")
	path := filepath.Join(modelsDir, modelName, version, "model.onnx")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("model not found: %s", path)
	}
	session, err := ort.NewDynamicAdvancdSession(
		path,
		[]string{"input"},
		[]string{"output"},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create onnx session:%w", err)
	}
	return session, nil
}
