package storage

import (
	"io/fs"
	"mime/multipart"
	"os"
)

type VideoStorage interface {
	StorageDetails() string
	GetVideo(path string) (*os.File, fs.FileInfo, error)
	UploadVideo(file multipart.File, file_ext string) (string, error)
	//DeleteVideo(filename string) error
}
