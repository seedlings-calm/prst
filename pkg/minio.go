package pkg_minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

type MinIOClient struct {
	Client   *minio.Client
	Endpoint string
	Bucket   string
	UseSSL   bool
}

// NewMinIOClient initializes a new MinIO client
func NewMinIOClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOClient{
		Client:   client,
		Endpoint: endpoint,
		Bucket:   bucket,
		UseSSL:   useSSL,
	}, nil
}

// EnsureBucket makes sure the bucket exists
func (m *MinIOClient) EnsureBucket(ctx context.Context) error {
	exists, err := m.Client.BucketExists(ctx, m.Bucket)
	if err != nil {
		return err
	}
	if !exists {
		err = m.Client.MakeBucket(ctx, m.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// UploadFile uploads a file to MinIO from an io.Reader
func (m *MinIOClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (minio.UploadInfo, error) {
	return m.Client.PutObject(ctx, m.Bucket, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})
}

// UploadMultipartFile 上传multipart.File
func (m *MinIOClient) UploadMultipartFile(ctx context.Context, objectName string, file multipart.File, size int64, contentType string) (minio.UploadInfo, error) {
	defer file.Close()
	return m.UploadFile(ctx, objectName, file, size, contentType)
}

// DownloadFile downloads an object as a stream
func (m *MinIOClient) DownloadFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	return m.Client.GetObject(ctx, m.Bucket, objectName, minio.GetObjectOptions{})
}

// GeneratePresignedURL creates a presigned URL for accessing a file
func (m *MinIOClient) GeneratePresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	u, err := m.Client.PresignedGetObject(ctx, m.Bucket, objectName, expiry, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// DeleteFile deletes a file from MinIO
func (m *MinIOClient) DeleteFile(ctx context.Context, objectName string) error {
	return m.Client.RemoveObject(ctx, m.Bucket, objectName, minio.RemoveObjectOptions{})
}

// ListFiles lists all objects in a bucket with a given prefix
func (m *MinIOClient) ListFiles(ctx context.Context, prefix string) []minio.ObjectInfo {
	var objects []minio.ObjectInfo
	for obj := range m.Client.ListObjects(ctx, m.Bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if obj.Err != nil {
			log.Println("list object error:", obj.Err)
			continue
		}
		objects = append(objects, obj)
	}
	return objects
}

// CheckFileExists checks if an object exists
func (m *MinIOClient) CheckFileExists(ctx context.Context, objectName string) (bool, error) {
	_, err := m.Client.StatObject(ctx, m.Bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetPublicURL returns the public URL of a file if the bucket is public
func (m *MinIOClient) GetPublicURL(objectName string) string {
	scheme := "http"
	if m.UseSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, m.Endpoint, m.Bucket, objectName)
}

// 设置对象的生命周期策略（删除超过指定时间的对象）
func (m *MinIOClient) SetBucketLifecycle(ctx context.Context) error {
	// 设置生命周期规则
	// Set lifecycle on a bucket
	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{
		{
			ID:     "bucket-lifecycle-rule",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: 365,
			},
		},
	}
	err := m.Client.SetBucketLifecycle(ctx, m.Bucket, config)
	if err != nil {
		return fmt.Errorf("设置生命周期时发生错误: %v", err)
	}
	return nil
}

// 模拟文件夹功能，通过使用前缀（Object Prefix）
func (m *MinIOClient) ListObjectsWithPrefix(bucketName, prefix string) ([]minio.ObjectInfo, error) {
	objectCh := m.Client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	var objects []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("列出对象时发生错误: %v", object.Err)
		}
		objects = append(objects, object)
	}
	return objects, nil
}
