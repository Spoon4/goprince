package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexRoute(t *testing.T) {
	router := initRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "This is the Prince RESTful API", w.Body.String())
}

func TestPrinceGenerateRoute(t *testing.T) {

	buffer := new(bytes.Buffer)
	mw := multipart.NewWriter(buffer)
	ioWriter, err := mw.CreateFormFile("html", "test.html")

	if assert.NoError(t, err) {
		ioWriter.Write([]byte("bin/test.html"))
	}
	ioWriter, err = mw.CreateFormFile("css", "test.css")

	if assert.NoError(t, err) {
		ioWriter.Write([]byte("bin/test.css"))
	}
	mw.Close()

	router := initRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/prince/generate/test.pdf", buffer)
	request.Header.Set("Content-Type", mw.FormDataContentType())

	router.ServeHTTP(w, request)

	assert.Equal(t, 200, w.Code)
}
