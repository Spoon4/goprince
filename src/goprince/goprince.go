package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os/signal"
	"syscall"
	"log"
)

const (
	TMP_DIR         = "/tmp/"
	DEFAULT_LOG_DIR = "/var/log/goprince/"
)

var (
	logDir string
	stdout bool
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

	wrapper := NewWrapper(htmlPath, logDir, stdout)
	license(&wrapper)

	cssPath, err := getFormFile(c, "stylesheet", true)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if cssPath != "" {
		wrapper.AddStyleSheet(cssPath)
	}

	//for multiple files upload
	/*
		cssFiles := c.MultipartForm().File["stylesheets[]"]
		if cssFiles != nil {
			for __, cssFile := range cssFiles {
				cssPath := filepath.Join(TMP_DIR, cssFile.Filename)
				c.SaveUploadedFile(cssFile, cssPath)
				wrapper.AddStyleSheet(cssPath)
			}
		}
	*/

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
// Set gin in release mode if APP_ENV var is set on 'production'
func main() {

	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show this message.")

	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var port int
	flag.StringVar(&logDir, "log-dir", DEFAULT_LOG_DIR, "Directory where log files must be stored.")
	flag.BoolVar(&stdout, "stdout", true, "If set, logs are displayed on stdout.")
	flag.IntVar(&port, "port", 8080, "Set Gin listening port")
	flag.Parse()

	f, _ := os.Create(filepath.Join(logDir, "gin.log"))

	if stdout {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	} else {
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(f)
	}

	env := os.Getenv("APP_ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// initialize Gin API routes
	router := initRouter()

	// manage signals to gracefully stop the application
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		os.Exit(0)
	}()

	// listen for HTTP requests
	log.Println(port)
	router.Run(fmt.Sprintf(":%d", port))
}
