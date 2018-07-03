package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"path/filepath"
)

func TestIndexRoute(t *testing.T) {
	router := initRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "This is a Prince RESTful API", w.Body.String())
}

func TestPrinceGenerateRoute(t *testing.T) {

	buffer := new(bytes.Buffer)
	mw := multipart.NewWriter(buffer)

	ioWriter, err := mw.CreateFormFile("input_file", "test.html")
	if assert.NoError(t, err) {
		ioWriter.Write([]byte("bin/test.html"))
	}

	ioWriter, err = mw.CreateFormFile("stylesheet", "test.css")
	if assert.NoError(t, err) {
		ioWriter.Write([]byte("bin/test.css"))
	}
	mw.Close()

	router := initRouter()

	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/generate/test.pdf", buffer)
	request.Header.Set("Content-Type", mw.FormDataContentType())
	router.ServeHTTP(w, request)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, filepath.Join(OUTPUT_DEST, "test.pdf"), w.Body.String())
	/*
		w2 := httptest.NewRecorder()
		request2, _ := http.NewRequest("POST", "/generate/test.pdf", buffer)
		request2.Header.Set("Content-Type", mw.FormDataContentType())
		request2.URL.Query().Add("output", "stream")
		router.ServeHTTP(w2, request2)
		assert.Equal(t, http.StatusOK, w2.Code)

		w3 := httptest.NewRecorder()
		request3, _ := http.NewRequest("POST", "/generate/test.pdf", buffer)
		request3.Header.Set("Content-Type", mw.FormDataContentType())
		request3.URL.Query().Add("output", "file")
		router.ServeHTTP(w3, request3)
		assert.Equal(t, http.StatusOK, w3.Code)

		//request.URL.Path = "/generate/"
	*/
}
