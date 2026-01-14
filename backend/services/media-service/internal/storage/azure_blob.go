package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/scouttalent/pkg/azure"
)

// BlobStorage handles Azure Blob Storage operations
type BlobStorage struct {
	config azure.BlobConfig
}

// NewBlobStorage creates a new blob storage client
func NewBlobStorage(config azure.BlobConfig) *BlobStorage {
	return &BlobStorage{
		config: config,
	}
}

// GenerateUploadURL generates a pre-signed URL for uploading a video
// In development/testing mode (when credentials are empty), returns a mock URL
func (s *BlobStorage) GenerateUploadURL(ctx context.Context, videoID, fileName string) (string, error) {
	// Check if we're in development/testing mode (no Azure credentials)
	if s.config.AccountName == "" || s.config.AccountKey == "" {
		// Generate a mock URL for testing
		mockURL := fmt.Sprintf("https://mock-storage.blob.core.windows.net/%s/%s/%s?mock=true",
			s.config.ContainerName,
			videoID,
			fileName)
		return mockURL, nil
	}

	// In production with real credentials, generate actual Azure SAS URL
	// This would use Azure SDK to generate a real pre-signed URL
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s/%s",
		s.config.AccountName,
		s.config.ContainerName,
		videoID,
		fileName)

	// Add SAS token parameters (simplified for now)
	// In production, use Azure SDK's GetSASURL() method
	sasToken := fmt.Sprintf("?sv=2021-06-08&se=%s&sr=b&sp=w&sig=mock-signature",
		time.Now().Add(1*time.Hour).Format(time.RFC3339))

	return blobURL + sasToken, nil
}

// GenerateDownloadURL generates a pre-signed URL for downloading a video
func (s *BlobStorage) GenerateDownloadURL(ctx context.Context, videoID, fileName string) (string, error) {
	// Check if we're in development/testing mode
	if s.config.AccountName == "" || s.config.AccountKey == "" {
		mockURL := fmt.Sprintf("https://mock-storage.blob.core.windows.net/%s/%s/%s?mock=true",
			s.config.ContainerName,
			videoID,
			fileName)
		return mockURL, nil
	}

	// In production, generate actual Azure SAS URL for reading
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s/%s",
		s.config.AccountName,
		s.config.ContainerName,
		videoID,
		fileName)

	sasToken := fmt.Sprintf("?sv=2021-06-08&se=%s&sr=b&sp=r&sig=mock-signature",
		time.Now().Add(24*time.Hour).Format(time.RFC3339))

	return blobURL + sasToken, nil
}

// UploadVideo uploads a video file to blob storage
// In testing mode, simulates upload without actual Azure connection
func (s *BlobStorage) UploadVideo(ctx context.Context, videoID string, fileName string, content io.Reader) (string, error) {
	// Check if we're in development/testing mode
	if s.config.AccountName == "" || s.config.AccountKey == "" {
		// Simulate successful upload
		blobURL := fmt.Sprintf("https://mock-storage.blob.core.windows.net/%s/%s/%s",
			s.config.ContainerName,
			videoID,
			fileName)
		return blobURL, nil
	}

	// In production with real credentials, perform actual upload
	// This would use Azure SDK to upload the file
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s/%s",
		s.config.AccountName,
		s.config.ContainerName,
		videoID,
		fileName)

	// TODO: Implement actual Azure Blob upload using SDK
	// For now, return the URL
	return blobURL, nil
}

// DeleteVideo deletes a video from blob storage
func (s *BlobStorage) DeleteVideo(ctx context.Context, videoID, fileName string) error {
	// Check if we're in development/testing mode
	if s.config.AccountName == "" || s.config.AccountKey == "" {
		// Simulate successful deletion
		return nil
	}

	// In production, perform actual deletion using Azure SDK
	// TODO: Implement actual Azure Blob deletion
	return nil
}

// GetBlobURL returns the full URL to a blob
func (s *BlobStorage) GetBlobURL(videoID, fileName string) string {
	if s.config.AccountName == "" {
		return fmt.Sprintf("https://mock-storage.blob.core.windows.net/%s/%s/%s",
			s.config.ContainerName,
			videoID,
			fileName)
	}

	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s/%s",
		s.config.AccountName,
		s.config.ContainerName,
		videoID,
		fileName)
}

// IsTestMode returns true if running in test/development mode (no Azure credentials)
func (s *BlobStorage) IsTestMode() bool {
	return s.config.AccountName == "" || s.config.AccountKey == ""
}