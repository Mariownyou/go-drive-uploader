package main

import (
	"fmt"
	"os"

	"github.com/mariownyou/go-drive-uploader/drive_uploader"
)

func main() {
	filename := "video.mp4"
	file, _ := os.ReadFile(filename)

	creds := []byte(os.Getenv("GOOGLE_DRRIVE_CREDENTIALS"))
	uploader, _ := drive_uploader.New(creds)

	link, _, _ := uploader.ShareFile(file, filename)
	fmt.Println(link)
}
