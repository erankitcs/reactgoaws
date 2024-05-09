package localstorage

import (
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	RootPath string
}

func (ls *LocalStorage) StorageDetails() string {
	return ls.RootPath
}

func (ls *LocalStorage) UploadVideo(file multipart.File, file_ext string) (string, error) {
	fmt.Println("Uploading video to local storage")
	// Generate a new UUID
	uuid := uuid.New()
	filename := strings.Replace(uuid.String(), "-", "", -1) + "." + file_ext
	fmt.Printf("Writing file with filename- %s", filename)
	dst := fmt.Sprintf("%s/%s", ls.RootPath, filename)
	destFile, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		return "", err
	}
	return filename, nil
}

// Write a function to return a video from storage filter by path
func (ls *LocalStorage) GetVideo(path string) (*os.File, fs.FileInfo, error) {
	fmt.Println("Getting video from local storage")
	videoPath := fmt.Sprintf("%s/%s", ls.RootPath, path)
	// Read a video object
	video, err := os.Open(videoPath)
	if err != nil {
		return nil, nil, err
	}
	// Get the file information
	videoInfo, err := video.Stat()
	if err != nil {
		return nil, nil, err
	}
	return video, videoInfo, nil
}
