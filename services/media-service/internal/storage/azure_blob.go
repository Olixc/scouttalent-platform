package storage

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/scouttalent/pkg/azure"
)

type AzureBlobClient struct {
	client        *azblob.Client
	containerName string
	accountName   string
}

func NewAzureBlobClient(cfg azure.BlobConfig) (*AzureBlobClient, error) {
	// Create credential
	credential, err := azblob.NewSharedKeyCredential(cfg.AccountName, cfg.AccountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	// Create service URL
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", cfg.AccountName)

	// Create client
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &AzureBlobClient{
		client:        client,
		containerName: cfg.ContainerName,
		accountName:   cfg.AccountName,
	}, nil
}

func (c *AzureBlobClient) GetBlobURL(blobName string) string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s",
		c.accountName, c.containerName, blobName)
}

func (c *AzureBlobClient) DeleteBlob(ctx context.Context, blobName string) error {
	_, err := c.client.DeleteBlob(ctx, c.containerName, blobName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}
	return nil
}

func (c *AzureBlobClient) BlobExists(ctx context.Context, blobName string) (bool, error) {
	_, err := c.client.NewBlobClient(c.containerName, blobName).GetProperties(ctx, nil)
	if err != nil {
		// Check if error is "blob not found"
		return false, nil
	}
	return true, nil
}