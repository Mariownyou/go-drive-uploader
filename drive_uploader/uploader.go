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
}

func New(credentialsJSON []byte) *Uploader {
	return &Uploader{
		credentialsJSON: credentialsJSON,
	}
}

func (u *Uploader) Upload(file []byte, filename string) string {
	ctx := context.Background()
	options := option.WithCredentialsJSON(u.credentialsJSON)
	service, err := drive.NewService(ctx, options)
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
	}

	// folderName := "My Folder"
	// folder := &drive.File{
	// 	Name:     folderName,
	// 	MimeType: "application/vnd.google-apps.folder",
	// }

	// createdFolder, err := service.Files.Create(folder).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to create folder: %v", err)
	// }

	fileReader := strings.NewReader(string(file))

	splitted := strings.Split(filename, ".")
	extension := splitted[len(splitted)-1]

	mimeType := mime.TypeByExtension("." + extension)

	driveFile := &drive.File{
		Name:     "My Media File.mp4",
		MimeType: mimeType,
		// Parents:  []string{createdFolder.Id},
	}

	uploadedFile, err := service.Files.Create(driveFile).Media(fileReader).Do()
	if err != nil {
		log.Fatalf("Unable to upload media file: %v", err)
	}

	permission := &drive.Permission{
		Type:               "anyone",
		Role:               "reader",
		AllowFileDiscovery: false,
	}
	_, err = service.Permissions.Create(uploadedFile.Id, permission).Do()
	if err != nil {
		log.Fatalf("Failed to create permission for the file: %v", err)
	}

	fileLink := fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", uploadedFile.Id)
	return fileLink
}
