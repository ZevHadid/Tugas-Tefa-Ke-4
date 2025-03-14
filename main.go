package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Read environment variables
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	// Inisialisasi client MinIO
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Contoh penggunaan fungsi-fungsi
	ctx := context.Background()

	// Upload file
	uploadFile(minioClient, ctx, "local-file.txt", "uploaded-file.txt")

	// Upload file dengan nama object yang mengandung /public
	uploadFile(minioClient, ctx, "local-file.txt", "public/uploaded-file.txt")

	// List file dalam bucket
	listFiles(minioClient, ctx)

	// Download file
	downloadFile(minioClient, ctx, "uploaded-file.txt", "downloaded-file.txt")

	// Generate presigned URL
	presignedURL := generatePresignedURL(minioClient, ctx, "uploaded-file.txt", 24*time.Hour)
	fmt.Println("Presigned URL:", presignedURL)

	// Generate public URL
	publicURL := generatePublicURL(endpoint, bucketName, "public/uploaded-file.txt")
	fmt.Println("Public URL:", publicURL)

	// Hapus file
	deleteFile(minioClient, ctx, "uploaded-file.txt")
}

func uploadFile(minioClient *minio.Client, ctx context.Context, filePath, objectName string) {
	_, err := minioClient.FPutObject(ctx, os.Getenv("MINIO_BUCKET_NAME"), objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("File %s berhasil diupload sebagai %s\n", filePath, objectName)
}

func listFiles(minioClient *minio.Client, ctx context.Context) {
	objectCh := minioClient.ListObjects(ctx, os.Getenv("MINIO_BUCKET_NAME"), minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}
		fmt.Println(object.Key)
	}
}

func downloadFile(minioClient *minio.Client, ctx context.Context, objectName, filePath string) {
	err := minioClient.FGetObject(ctx, os.Getenv("MINIO_BUCKET_NAME"), objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("File %s berhasil didownload sebagai %s\n", objectName, filePath)
}

func generatePresignedURL(minioClient *minio.Client, ctx context.Context, objectName string, expiry time.Duration) string {
	presignedURL, err := minioClient.PresignedGetObject(ctx, os.Getenv("MINIO_BUCKET_NAME"), objectName, expiry, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return presignedURL.String()
}

func generatePublicURL(endpoint, bucketName, objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, objectName)
}

func deleteFile(minioClient *minio.Client, ctx context.Context, objectName string) {
	err := minioClient.RemoveObject(ctx, os.Getenv("MINIO_BUCKET_NAME"), objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("File %s berhasil dihapus\n", objectName)
}
