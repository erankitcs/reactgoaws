package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type JSONReponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 // one MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must contain a single JSON value")
	}
	return nil
}

func (app *application) readMultiPartForm(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	maxMultiPartBytes := 1024 * 1024 * 20 // 20 MB
	r.ParseMultipartForm(int64(maxMultiPartBytes))
	// Get the uploaded file from the request
	file, fileheader, err := r.FormFile("movievideofile")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	return file, fileheader, nil

}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONReponse
	payload.Error = true
	payload.Message = err.Error()
	fmt.Println(err)
	return app.writeJSON(w, statusCode, payload)
}
