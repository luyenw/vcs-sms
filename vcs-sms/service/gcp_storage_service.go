package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type ClientUploader struct {
	cli        *storage.Client
	bucketName string
	uploadPath string
}

func (c *ClientUploader) UploadFileAndSetMetaData(file *os.File, object string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := c.cli.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	// Update the object to set the metadata.
	o := c.cli.Bucket(c.bucketName).Object(c.uploadPath + object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		ACL: []storage.ACLRule{
			{Entity: storage.AllUsers, Role: storage.RoleReader},
		},
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return fmt.Errorf("ObjectHandle(%q).Update: %w", object, err)
	}
	return nil
}

func (c *ClientUploader) GetFileURL(object string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := c.cli.Bucket(c.bucketName).Object(c.uploadPath + object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("Object(%q).Attrs: %w", object, err)
	}
	return attrs.MediaLink, nil
}
