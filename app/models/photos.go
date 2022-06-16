package models

import (
	"encoding/base64"
	"net/http"
)

type Photos []string

// GetFileTypes processing files is slow/ expensive, I typically wouldn't do it in an HTTP API request handler. Instead,
// either stick the encoded files in a database, or pass to another process via some channel directly for asynchronous processing.
func (p *Photos) GetFileTypes() ([]string, error) {
	photos := make([]string, 0)
	for _, photo := range *p {
		decodedPhoto, err := base64.StdEncoding.DecodeString(photo)
		if err != nil {
			return photos, err
		}

		// not the best way to figure out the MIME type, but the only way that comes to mind without
		// relying on 3rd party libraries with C underpinnings, written by people much smarter than me!
		photos = append(photos, http.DetectContentType(decodedPhoto))
	}
	return photos, nil
}
