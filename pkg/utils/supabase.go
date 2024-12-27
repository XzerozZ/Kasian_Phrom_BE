package utils

import (
	"fmt"
	"os"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"

	storage_go "github.com/supabase-community/storage-go"
)

func UploadImage(fileName string, dir string, config configs.Supabase) (string, error) {
    filePath := "./uploads/" + fileName
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to open file: %w", err)
    }

    defer file.Close()
    if config.URL == "" || config.Key == "" || config.Bucket == "" {
        return "", fmt.Errorf("invalid Supabase config: URL='%s', Key='%s', Bucket='%s'", config.URL, config.Key, config.Bucket)
    }
    
    storageClient := storage_go.NewClient(config.URL, config.Key, nil)
    if storageClient == nil {
        return "", fmt.Errorf("failed to create storage client: invalid Supabase configuration")
    }

	options := storage_go.FileOptions{
        ContentType: func() *string { s := "image/jpeg"; return &s }(),
    }

    fileName = dir + "/" + fileName
    _, err = storageClient.UploadFile(config.Bucket, fileName, file, options)
    if err != nil {
        return "", fmt.Errorf("failed to upload file '%s' to bucket '%s': %w", fileName, config.Bucket, err)
    }

    url := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", config.URL, config.Bucket, fileName)
    return url, nil
}