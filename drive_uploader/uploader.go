package drive_uploader

import (
	"context"
	"fmt"
	"log"
	"mime"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Uploader struct {
	credentialsJSON []byte
	service         *drive.Service
}

func New(credentialsJSON []byte) (*Uploader, error) {
	ctx := context.Background()
	options := option.WithCredentialsJSON(credentialsJSON)
	service, err := drive.NewService(ctx, options)

	if err != nil {
		log.Fatalf("Unable to create Drivade service: %v", err)
		return nil, err
	}

	return &Uploader{
		credentialsJSON: credentialsJSON,
		service:         service,
	}, nil
}

func (u *Uploader) Upload(file []byte, filename string) (string, string, error) {
	fileReader := strings.NewReader(string(file))

	splitted := strings.Split(filename, ".")
	extension := splitted[len(splitted)-1]

	mimeType := mime.TypeByExtension("." + extension)

	driveFile := &drive.File{
		Name:     filename,
		MimeType: mimeType,
	}

	uploadedFile, err := u.service.Files.Create(driveFile).Media(fileReader).Do()
	if err != nil {
		log.Fatalf("Unable to upload media file: %v", err)
		return "", "", err
	}

	permission := &drive.Permission{
		Type:               "anyone",
		Role:               "reader",
		AllowFileDiscovery: false,
	}
	_, err = u.service.Permissions.Create(uploadedFile.Id, permission).Do()
	if err != nil {
		log.Fatalf("Failed to create permission for the file: %v", err)
		return "", "", err
	}

	fileLink := fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", uploadedFile.Id)
	return fileLink, uploadedFile.Id, nil
}

func (u *Uploader) Delete(fileID string) error {
	err := u.service.Files.Delete(fileID).Do()
	if err != nil {
		return err
	}

	return nil
}
