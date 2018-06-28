package main

import (
	//"fmt"
	// "strings"
	// "log"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

const TMP_DIR = "/tmp/"

// Default index handler
func indexHandler(c *gin.Context) {
	c.String(http.StatusOK, "This is a Prince RESTful API")
}

// API handler to generate PDF from HTML files
func generateHandler(c *gin.Context) {

	outputFile := c.Param("filename")

	// Get files from POST data and save them in temp dir
	htmlFile, _ := c.FormFile("html")
	htmlPath := filepath.Join(TMP_DIR, htmlFile.Filename)
	c.SaveUploadedFile(htmlFile, htmlPath)

	cssFile, _ := c.FormFile("css")
	cssPath := filepath.Join(TMP_DIR, cssFile.Filename)
	c.SaveUploadedFile(cssFile, cssPath)

	wrapper := NewWrapper(htmlPath)
	license(&wrapper)

	//cssFiles := c.MultipartForm().File["css[]"]
	//if cssFiles != nil {
	//	for __, cssFile := range cssFile {
	//		cssPath := filepath.Join(TMP_DIR, cssFile.Filename)
	//		c.SaveUploadedFile(cssFile, cssPath)
	wrapper.AddStyleSheet(cssPath)
	//	}
	//}

	dest := wrapper.Generate(outputFile)

	// Remove temp files
	os.Remove(htmlPath)
	os.Remove(cssPath)

	//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+outputFile)
	//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", outputFile))
	c.Header("Content-Type", "application/pdf")
	//c.Header("Content-Type", "application/octet-stream")
	c.File(dest)
}

// Manage Prince license
// Prince key or file comes from env vars
// LICENSE_KEY / LICENSE_FILE
func license(wrapper *Wrapper) {
	licenseFile := os.Getenv("LICENSE_FILE")
	licenseKey := os.Getenv("LICENSE_KEY")

	if "" != licenseFile {
		(*wrapper).SetLicenseFile(licenseFile)
	}
	if "" != licenseKey {
		(*wrapper).SetLicenseKey(licenseKey)
	}
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
	}

	router := initRouter()
	router.Run()
}
