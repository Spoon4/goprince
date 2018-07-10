package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"os/signal"
	"syscall"
	"io/ioutil"

	"github.com/goprince/src/goprince/storage"
)

const (
	TMP_DIR         = "/tmp/"
	DEFAULT_LOG_DIR = "/var/log/goprince/"
	DEFAULT_STORAGE = STORAGE_MODE_LOCAL
)

const (
	STORAGE_MODE_LOCAL     = "local"
	STORAGE_MODE_AZURE     = "azure"
	STORAGE_MODE_OPENSTACK = "openstack"
)

// shared vars
var (
	logDir  string
	stdout  bool
	store storage.PrinceStorage
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

func changeStorage(mode string) {

	switch mode {
	case STORAGE_MODE_AZURE:
		store = storage.NewAzureClient()
		log.Println("storeage mode: Azure")
		break
	case STORAGE_MODE_OPENSTACK:
		store = storage.NewOpenstackClient()
		log.Println("storeage mode: Openstack")
		break
	default:
		log.Println("storeage mode: Local")
	}
}

func writeToContainer(c *gin.Context) {
	path, _ := getFormFile(c, "file", false)
	filename := c.Param("filename")
	container := c.Param("application")

	file, err := ioutil.ReadFile(path)
	if nil != err {
		c.String(http.StatusInternalServerError, err.Error())
	}
	store.Write(filename, container, file)
	os.Remove(path)
	c.String(http.StatusCreated, filename)
}

// Gin router initialization.
// Bind all API routes with their handler
func initRouter() *gin.Engine {

	router := gin.Default()
	router.GET("/", indexHandler)
	router.POST("/generate/:filename", generateHandler)
	router.POST("/write/:application/:filename", writeToContainer)
	return router
}

// Main entrypoint
// Manage command line flags and help message
// Set up Gin in release mode if APP_ENV var is set on 'production'
func main() {

	// display help if needed
	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show this message.")

	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// get command line args/options
	var port int
	var storageMode string
	flag.StringVar(&logDir, "log-dir", DEFAULT_LOG_DIR, "Directory where log files must be stored.")
	flag.StringVar(&storageMode, "storage", DEFAULT_STORAGE, "Set storage mode")
	flag.BoolVar(&stdout, "stdout", true, "If set, logs are displayed on stdout.")
	flag.IntVar(&port, "port", 8080, "Set Gin listening port")
	flag.Parse()

	// configure Gin logging and create log file
	f, _ := os.Create(filepath.Join(logDir, "gin.log"))

	if stdout {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	} else {
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(f)
	}

	// set application environement from env var
	env := os.Getenv("APP_ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	changeStorage(storageMode)

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
	log.Println(fmt.Sprintf("Goprince is listening on port %d...", port))
	router.Run(fmt.Sprintf(":%d", port))
}
