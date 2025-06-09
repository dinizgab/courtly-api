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
	Client *supabase.Client
	Bucket string
}

func NewSupabaseStorageUploader(config config.StorageConfig, bucket string) StorageUploader {
	client := supabase.NewClient(
		config.ProjectURL,
		config.APIKey,
		nil,
	)

	return &supabaseStorageUploader{
		Client: client,
		Bucket: bucket,
	}
}

func (s *supabaseStorageUploader) UploadFile(ctx context.Context, courtId string, filename string, fileBytes io.Reader) (string, error) {
	companyId := ctx.Value("company_id").(string)

	remoteFilePath := fmt.Sprintf("%s/%s/%s", companyId, courtId, filename)

	_, err := s.Client.UploadFile(s.Bucket, companyId, fileBytes)
	if err != nil {
		return "", fmt.Errorf("SupabaseStorageUploader.UploadFile - Error uploading file: %w", err)
	}

	url := s.Client.GetPublicUrl(s.Bucket, remoteFilePath)

	return url.SignedURL, nil
}
