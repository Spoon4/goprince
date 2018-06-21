package main

import (
	// "fmt"
	// "strings"
	// "log"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const TMP_DIR = "/tmp/"
const PRINCE_BIN = "prince"

func indexHandler(context *gin.Context) {
	context.String(http.StatusOK, "This is the Prince RESTful API")
}

func generateHandler(context *gin.Context) {

	_, err := exec.LookPath(PRINCE_BIN)
	if err != nil {
		context.String(http.StatusInternalServerError, "didn't find '%s' executable", PRINCE_BIN)
		return

	}

	// Get files from POST data and save them in temp dir
	htmlFile, _ := context.FormFile("html")
	cssFile, _ := context.FormFile("css")

	htmlPath := filepath.Join(TMP_DIR, htmlFile.Filename)
	cssPath := filepath.Join(TMP_DIR, cssFile.Filename)
	outputFile := filepath.Join("/public/", context.Params.ByName("filename"))

	context.SaveUploadedFile(htmlFile, htmlPath)
	context.SaveUploadedFile(cssFile, cssPath)

	// Execute Prince command
	_, err = exec.Command(PRINCE_BIN, "-s", cssPath, htmlPath, "-o", outputFile).Output()

	// Remove temp files
	os.Remove(htmlPath)
	os.Remove(cssPath)

	if err != nil {
		context.String(http.StatusInternalServerError, err.Error()+"\n")
		return
	}

	//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Content-Disposition", "attachment; filename="+outputFile)
	context.Header("Content-Type", "application/pdf")
	context.File(outputFile)
}

func initRouter() *gin.Engine {

	router := gin.Default()
	router.GET("/", indexHandler)

	prince := router.Group("/prince")
	{
		prince.POST("/generate/:filename", generateHandler)
	}

	return router
}

func main() {

	env := os.Getenv("APP_ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := initRouter()
	router.Run()
}
