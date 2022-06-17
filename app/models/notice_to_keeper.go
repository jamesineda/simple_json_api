package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

var ErrInvalidFileFormat = fmt.Errorf("invalid file format - must be PDF")
var ErrMustSupplyEitherFileOrURL = fmt.Errorf("neither file or url have been supplied")

type NoticeToKeeper struct {
	File string `json:"file"`
	Url  string `json:"url"`
}

func (ntk *NoticeToKeeper) UnmarshalJSON(data []byte) error {
	type Alias NoticeToKeeper

	aux := struct {
		File string `json:"file"`
		*Alias
	}{
		Alias: (*Alias)(ntk),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.File != "" {
		decodedFile, err := base64.StdEncoding.DecodeString(aux.File)
		if err != nil {
			return err
		}

		// I would be interested to know if there is a better way to do this? :)
		if http.DetectContentType(decodedFile) != "application/pdf" {
			return ErrInvalidFileFormat
		}
	}

	if aux.Url == "" && aux.File == "" {
		return ErrMustSupplyEitherFileOrURL
	}

	ntk.File = aux.File

	return nil
}
