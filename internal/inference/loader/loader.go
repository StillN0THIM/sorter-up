package loader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
	ort "github.com/yalue/onnxruntime_go"
)

func sessionKey(modelName, version string) string {
	return fmt.Sprintf("session:%s:%s", modelName, version)
}
