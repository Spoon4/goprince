package main

import (
	"fmt"
	// "strings"
	// "log"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"io"
)

const (
	TMP_DIR = "/tmp/"
	LOG_DIR = "/var/log/goprince/"
)

// Default index handler
func indexHandler(c *gin.Context) {
	c.String(http.StatusOK, "This is a Prince RESTful API")
}

// API handler to generate PDF from HTML files
func generateHandler(c *gin.Context) {

	outputFile := c.Param("filename")

	htmlPath, err := getFormFile(c, "input_file", false)

	if err != nil {
		c.String(http.StatusInternalServerError, "No file to convert")
		return
	}

	wrapper := NewWrapper(htmlPath, LOG_DIR)
	license(&wrapper)

	cssPath, err := getFormFile(c, "stylesheet", true)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if cssPath != "" {
		wrapper.AddStyleSheet(cssPath)

		//for multiple files upload
		/*
			cssFiles := c.MultipartForm().File["stylesheets[]"]
			if cssFiles != nil {
				for __, cssFile := range cssFile {
					cssPath := filepath.Join(TMP_DIR, cssFile.Filename)
					c.SaveUploadedFile(cssFile, cssPath)
					wrapper.AddStyleSheet(cssPath)
				}
			}
		*/
	}

	dest, err := wrapper.Generate(outputFile)

	if nil != err {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Remove temp files
	os.Remove(htmlPath)
	os.Remove(cssPath)

	output := c.DefaultQuery("output", "")
	switch output {
	case "file":
		//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+outputFile)
		//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", outputFile))
		c.Header("Content-Type", "application/pdf")
		c.File(dest)
		break
	case "stream":
		//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Type", "application/octet-stream")
		c.File(dest)
		break
	default:
		c.String(http.StatusCreated, dest)
	}
}

// Manage Prince license
// Prince key or file comes from env vars
// LICENSE_KEY / LICENSE_FILE
func license(wrapper *Wrapper) {

	licenseFile := os.Getenv("LICENSE_FILE")
	licenseKey := os.Getenv("LICENSE_KEY")

	(*wrapper).SetLicenseFile(licenseFile)
	(*wrapper).SetLicenseKey(licenseKey)
}

// Get files from POST data and save them in temp dir
func getFormFile(c *gin.Context, parameter string, optional bool) (path string, err error) {

	file, _ := c.FormFile(parameter)

	if nil != file {
		path := filepath.Join(TMP_DIR, file.Filename)
		c.SaveUploadedFile(file, path)
		return path, nil
	} else {
		if false == optional {
			return "", errors.New(fmt.Sprintf("Parameter %s is not present in the form", parameter))
		}
	}
	return "", nil
}

// Gin router initialization.
// Bind all API routes with their handler
func initRouter() *gin.Engine {

	router := gin.Default()
	router.GET("/", indexHandler)
	router.POST("/generate/:filename", generateHandler)
	return router
}

// Main entrypoint
// Set gin release mode if APP_ENV var is set on 'production'
func main() {

	env := os.Getenv("APP_ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
	}

	f, _ := os.Create(filepath.Join(LOG_DIR, "gin.log"))
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := initRouter()
	router.Run()
}
