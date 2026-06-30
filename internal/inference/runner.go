package inference

import (
	"context"
	"fmt"
	"io"

	"github.com/StillN0THIM/sorter-up/internal/inference/loader"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	ort "github.com/yalue/onnxruntime_go"
)

func Run(ctx context.Context, rdb *redis.Client, modelName, version string, r io.Reader) (interface{}, error) {
	session, err := loader.GetSession(ctx, rdb, modelName, version)
	if err != nil {
		return nil, fmt.Errorf("load model: %w", err)
	}
	imageBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read input: %w", err)
	}

	inputTensor, err := preprocess(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("preprocess:%w", err)
	}
	outputs, err := session.Run([]ort.ArbitraryTensor{inputTensor})
	if err != nil {
		return nil, fmt.Errorf("onnx run: %w", err)
	}
	return postprocess(outputs)
}

func preprocess(imageBytes []byte) (*ort.Tensor[float32], error) {
	data := make([]float32, 1*3*224*224)
	shape := ort.NewShape(1, 3, 224, 224)
	return ort.NewTesnor(shape, data)
}

func postprocess(output []ort.ArbitraryTeanor) (interface{}, error) {
	out, ok := output[0].(*ort.Tensor[float32])
	if !ok {
		return nil, fmt.Errorf("unexpexted output tensor type")
	}
	return gin.H{
		"scores": out.GetData(),
	}, nil
}
