package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"osm-mngt/internal/config"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Init 初始化MinIO客户端
// 容错处理：连接失败时返回nil，不阻断服务启动
func Init(cfg config.MinIOConfig) (*minio.Client, error) {
	// 如果endpoint为空，直接返回nil
	if cfg.Endpoint == "" {
		log.Println("Warning: MinIO endpoint not configured, skipping MinIO initialization")
		return nil, nil
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Printf("Warning: Failed to create MinIO client: %v", err)
		return nil, nil
	}

	// 检查连接是否正常
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5秒超时
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		log.Printf("Warning: Failed to connect to MinIO: %v (service will continue without file upload)", err)
		return nil, nil
	}

	// 如果bucket不存在，尝试创建
	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Warning: Failed to create bucket %s: %v", cfg.BucketName, err)
			return nil, nil
		}
		log.Printf("Info: Created bucket %s successfully", cfg.BucketName)
	}

	log.Println("Info: MinIO initialized successfully")
	return client, nil
}

// UploadFile 上传文件到MinIO（带容错）
func UploadFile(client *minio.Client, bucketName string, ctx context.Context, objectName string, reader interface{}, objectSize int64, opts minio.PutObjectOptions) (string, error) {
	if client == nil {
		return "", fmt.Errorf("MinIO client not available")
	}

	_, err := client.PutObject(ctx, bucketName, objectName, reader.(io.Reader), objectSize, opts)
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	return objectName, nil
}
