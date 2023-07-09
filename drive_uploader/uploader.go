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

func (u *Uploader) ShareFile(file []byte, filename string) (string, string, error) {
	splitted := strings.Split(filename, ".")
	extension := splitted[len(splitted)-1]

	mimeType := mime.TypeByExtension("." + extension)

	f := &drive.File{
		Name:     filename,
		MimeType: mimeType,
	}

	p := &drive.Permission{
		Type:               "anyone",
		Role:               "reader",
		AllowFileDiscovery: false,
	}

	return u.Upload(file, f, p)
}

func (u *Uploader) Upload(b []byte, f *drive.File, p *drive.Permission) (string, string, error) {
	fileReader := strings.NewReader(string(b))

	uploadedFile, err := u.service.Files.Create(f).Media(fileReader).Do()
	if err != nil {
		log.Fatalf("Unable to upload media file: %v", err)
		return "", "", err
	}

	if p != nil {
		_, err = u.service.Permissions.Create(uploadedFile.Id, p).Do()
		if err != nil {
			log.Fatalf("Failed to create permission for the file: %v", err)
			return "", "", err
		}
	}

	fileLink := fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", uploadedFile.Id)
	return fileLink, uploadedFile.Id, nil
}

func (u *Uploader) CreateFolder(name string, parents ...string) (string, error) {

	f := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  parents,
	}

	createdFolder, err := u.service.Files.Create(f).Do()
	if err != nil {
		log.Fatalf("Unable to create folder: %v", err)
		return "", err
	}

	return createdFolder.Id, nil
}

func (u *Uploader) Delete(fileID string) error {
	err := u.service.Files.Delete(fileID).Do()
	if err != nil {
		return err
	}

	return nil
}
