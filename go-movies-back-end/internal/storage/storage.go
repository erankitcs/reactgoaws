package storage

import (
	"mime/multipart"
	"net/http"
	"os"
)

type VideoStorage interface {
	StorageDetails() string
	GetVideo(path string, w http.ResponseWriter) (os.FileInfo, error)
	UploadVideo(file multipart.File, file_ext string) (string, error)
	DeleteVideo(path string) error
}
