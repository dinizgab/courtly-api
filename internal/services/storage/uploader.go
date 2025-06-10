package storage

import (
	"context"
	"fmt"
	"github.com/dinizgab/booking-mvp/internal/config"
	supabase "github.com/supabase-community/storage-go"
	"io"
)

type StorageUploader interface {
	UploadFile(ctx context.Context, courtId string, filename string, fileBytes io.Reader) (string, error)
}

type supabaseStorageUploader struct {
	ProjectURL string
	APIKey     string
	Bucket     string
}

func NewSupabaseStorageUploader(config *config.StorageConfig, bucket string) StorageUploader {
	return &supabaseStorageUploader{
		ProjectURL: config.ProjectURL,
		APIKey:     config.APIKey,
		Bucket:     bucket,
	}
}

func (s *supabaseStorageUploader) UploadFile(ctx context.Context, courtId string, filename string, fileBytes io.Reader) (string, error) {
	token := ctx.Value("jwt_token").(string)
	companyId := ctx.Value("company_id").(string)
	remoteFilePath := fmt.Sprintf("%s/%s/%s", companyId, courtId, filename)
	client := supabase.NewClient(
		s.ProjectURL,
		s.APIKey,
		map[string]string{
			"Authorization": "Bearer " + token,
		},
	)

	_, err := client.UploadFile(s.Bucket, remoteFilePath, fileBytes)
	if err != nil {
		return "", fmt.Errorf("SupabaseStorageUploader.UploadFile - Error uploading file: %w", err)
	}

	url := client.GetPublicUrl(s.Bucket, remoteFilePath)

	return url.SignedURL, nil
}
