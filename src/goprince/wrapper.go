package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	PRINCE_BIN  = "prince"
	OUTPUT_DEST = "/public/"
	LOG_FILE    = "prince.log"
)

// Exporting interface instead of struct
type Wrapper interface {
	SetHTML(isHTML bool)
	AddStyleSheet(cssPath string)
	AddScript(scriptPath string)
	AddFileAttachment(file string)

	SetLicenseFile(file string)
	SetLicenseKey(hash string)

	SetLogFile(logPath string)

	SetPageSize(size string)
	SetPageMargin(margin string)

	SetPDFOutputIntent(profile string)
	SetPDFProfile(profile string)
	SetPDFTitle(title string)
	SetPDFSubject(subject string)
	SetPDFAuthor(author string)
	SetPDFKeywords(keywords string)
	SetPDFCreator(creator string)

	Generate(outputFile string) (outputPath string, err error)
}

// Struct are not exported
type Remap struct {
	url string
	dir string
}

// Define struct to configure and run Prince command line
type Prince struct {
	exePath   string
	inputFile string
	inputType string

	pageSize   string
	pageMargin string

	styleSheets     []string
	scripts         []string
	fileAttachments []string
	remaps          []Remap

	javascript bool
	insecure   bool

	licenseFile string
	licenseKey  string

	pdfOutputIntent string
	pdfProfile      string
	pdfTitle        string
	pdfSubject      string
	pdfAuthor       string
	pdfKeywords     string
	pdfCreator      string

	logFile string
	debug   bool
	verbose bool
	stdout  bool

	noWarnCssUnknown     bool
	noWarnCssUnsupported bool
}

func NewWrapper(inputFile string, logPath string, stdout bool) Wrapper {

	w := new(Prince)
	w.exePath = PRINCE_BIN
	w.logFile = filepath.Join(logPath, LOG_FILE)

	w.inputFile = inputFile
	w.inputType = "auto"

	w.styleSheets = make([]string, 0, 50)
	w.scripts = make([]string, 0, 50)
	w.fileAttachments = make([]string, 0, 50)
	w.remaps = make([]Remap, 0, 50)

	isDev := os.Getenv("APP_ENV") != "production"
	w.debug = isDev
	w.verbose = isDev
	w.stdout = stdout

	w.noWarnCssUnknown = !isDev
	w.noWarnCssUnsupported = !isDev
	return w
}

func (w *Prince) Generate(outputFile string) (outputPath string, err error) {

	_, err = exec.LookPath(w.exePath)
	if err != nil {
		log.Println(err.Error())
		return "", err

	}

	outputPath = filepath.Join(OUTPUT_DEST, outputFile)

	args := w.GetCommandLineArgs(outputPath)
	args = append(args, w.inputFile)

	logWriter, err := w.SetLogger()

	if nil != err {
		log.Println(err.Error())
		return "", err
	}

	cmd := exec.Command(w.exePath, args...)

	if w.stdout {
		cmd.Stdout = logWriter
		cmd.Stderr = logWriter
	}

	err = cmd.Run()

	if nil != err {
		log.Println(err.Error())
		return "", err
	}

	return outputPath, nil
}

// Set log file and create writer for command line output
// If stdout property is set to true, output is also redirected to stdout

func (w *Prince) SetLogger() (writer io.Writer, err error) {

	f, err := os.OpenFile(w.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening file: %v", err))
	}

	if w.stdout {
		writer = io.MultiWriter(f, os.Stdout)
	} else {
		writer = io.Writer(f)
	}
	return writer, nil
}

// Compute Prince command line arguments in a single array
// outputFile: filename of the output file
// Return the array of arguments
func (w *Prince) GetCommandLineArgs(outputFile string) []string {

	args := make([]string, 0)

	for _, stylesheet := range w.styleSheets {
		args = append(args, "-s", strconv.Quote(stylesheet))
	}
	for _, script := range w.scripts {
		args = append(args, "--script", strconv.Quote(script))
	}
	for _, attachment := range w.fileAttachments {
		args = append(args, "--attach="+strconv.Quote(attachment))
	}
	for _, remap := range w.remaps {
		args = append(args, "--remap="+strconv.Quote(remap.url)+"="+strconv.Quote(remap.dir))
	}
	if "auto" != w.inputType {
		args = append(args, "-i", w.inputType)
	}
	if true == w.javascript {
		args = append(args, "--javascript")
	}
	if true == w.insecure {
		args = append(args, "--insecure")
	}

	if "" != w.licenseKey {
		args = append(args, "--license-key="+strconv.Quote(w.licenseKey))
	}
	if "" != w.licenseFile {
		args = append(args, "--license-file="+strconv.Quote(w.licenseFile))
	}

	if "" != w.pdfOutputIntent {
		args = append(args, "--pdf-output-intent="+strconv.Quote(w.pdfOutputIntent))
	}
	if "" != w.pdfProfile {
		args = append(args, "--pdf-profile="+strconv.Quote(w.pdfProfile))
	}
	if "" != w.pdfTitle {
		args = append(args, "--pdf-title="+strconv.Quote(w.pdfTitle))
	}
	if "" != w.pdfSubject {
		args = append(args, "--pdf-subject="+strconv.Quote(w.pdfSubject))
	}
	if "" != w.pdfAuthor {
		args = append(args, "--pdf-author="+strconv.Quote(w.pdfAuthor))
	}
	if "" != w.pdfKeywords {
		args = append(args, "--pdf-keywords="+strconv.Quote(w.pdfKeywords))
	}
	if "" != w.pdfCreator {
		args = append(args, "--pdf-creator="+strconv.Quote(w.pdfCreator))
	}

	//if "" != w.logFile {
	//	args = append(args, "--log="+w.logFile)
	//}
	if true == w.debug {
		args = append(args, "--debug")
	}
	if true != w.verbose {
		args = append(args, "--verbose")
	}

	if true != w.noWarnCssUnknown {
		args = append(args, "--no-warn-css-unknown")
	}
	if true != w.noWarnCssUnsupported {
		args = append(args, "--no-warn-css-unsupported")
	}

	args = append(args, "--structured-log=normal", "-o", outputFile)

	return args
}

// License management by file
func (w *Prince) SetLicenseFile(file string) {
	w.licenseFile = file
}

// License management by hash key
func (w *Prince) SetLicenseKey(hash string) {
	w.licenseKey = hash
}

// Add new CSS file to embed in PDF
func (w *Prince) AddStyleSheet(cssPath string) {
	w.styleSheets = append(w.styleSheets, cssPath)
}

// Empty all CSS files embedded
func (w *Prince) ClearStyleSheets() {
}

// Add new javascript file to embed in PDF
func (w *Prince) AddScript(scriptPath string) {
	w.scripts = append(w.scripts, scriptPath)
}

// Empty all javascript files embedded
func (w *Prince) ClearScripts() {
}

// Add new attachment file with PDF
func (w *Prince) AddFileAttachment(file string) {
	w.fileAttachments = append(w.fileAttachments, file)
}

// Empty all attached files
func (w *Prince) ClearFileAttachments() {
}

// Define file input
func (w *Prince) SetHTML(isHTML bool) {
	if isHTML == true {
		w.inputType = "html"
	} else {
		w.inputType = "xml"
	}
}

func (w *Prince) SetPageSize(size string) {
	w.pageSize = size
}

func (w *Prince) SetPageMargin(margin string) {
	w.pageMargin = margin
}

func (w *Prince) SetPDFOutputIntent(profile string) {
	w.pdfOutputIntent = profile
}

func (w *Prince) SetPDFProfile(profile string) {
	w.pdfProfile = profile
}

func (w *Prince) SetPDFTitle(title string) {
	w.pdfTitle = title
}

func (w *Prince) SetPDFSubject(subject string) {
	w.pdfSubject = subject
}

func (w *Prince) SetPDFAuthor(author string) {
	w.pdfAuthor = author
}

func (w *Prince) SetPDFKeywords(keywords string) {
	w.pdfKeywords = keywords
}

func (w *Prince) SetPDFCreator(creator string) {
	w.pdfCreator = creator
}

// Specify a file that Prince should use to log error/warning messages.
// logFile: The filename that Prince should use to log error/warning
// messages, or '' to disable logging.
func (w *Prince) SetLogFile(logPath string) {
	w.logFile = logPath
}
