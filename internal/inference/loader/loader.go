package loader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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

type cache struct {
	mu    sync.RWMutex
	store map[string]*ort.DynamicAdvancedSession
}

var sessionCache = &cache{store: make(map[string]*ort.DynamicAdvancedSession)}

func (c *cache) get(key string) (*ort.DynamicAdvancedSession, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.store[key]
	return s, ok
}

func (c *cache) set(key string, s *ort.DynamicAdvancesSession) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = s
}
