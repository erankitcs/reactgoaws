package storage

import "mime/multipart"

type VideoStorage interface {
	StorageDetails() string
	//GetVideo(filename string) (string, error)
	UploadVideo(file multipart.File, file_ext string) (string, error)
	//DeleteVideo(filename string) error
}
