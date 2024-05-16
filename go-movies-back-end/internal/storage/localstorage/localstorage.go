package localstorage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	//dst := fmt.Sprintf("%s/%s", ls.RootPath, filename)
	// Sanitize the path to prevent directory traversal
	videoPath := filepath.FromSlash(filepath.Join(ls.RootPath, filename))

	destFile, err := os.Create(videoPath)
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
func (ls *LocalStorage) GetVideo(path string, w http.ResponseWriter) (os.FileInfo, error) {
	fmt.Println("Getting video from local storage")

	// Sanitize the path to prevent directory traversal
	videoPath := filepath.FromSlash(filepath.Join(ls.RootPath, path))

	// Read a video object
	video, err := os.Open(videoPath)
	if err != nil {
		return nil, err
	}
	defer video.Close()
	// Get the file information
	videoInfo, err := video.Stat()
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(w, video)
	if err != nil {
		return nil, err
	}
	return videoInfo, nil
}

// A fuction which will delete a video from storage
func (ls *LocalStorage) DeleteVideo(path string) error {
	fmt.Println("Deleting video from local storage")

	// Sanitize the path to prevent directory traversal
	videoPath := filepath.FromSlash(filepath.Join(ls.RootPath, path))

	// Remove a video object
	err := os.Remove(videoPath)
	if err != nil {
		return err
	}
	return nil
}
