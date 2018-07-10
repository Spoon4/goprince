package storage

import (
	"net/http"
)

// Exporting interface instead of struct
type PrinceStorage interface {
	GetConfig(filename string)
	Authenticate(username string, password string)
	Read(object string) []byte
	Write(filename string, container string, fContent []byte)
	Delete(object string)
	Metadata(object string) (metadata http.Header, err error)
	Debug(activate bool)
}
