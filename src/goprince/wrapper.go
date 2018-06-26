package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	PRINCE_BIN  = "prince"
	OUTPUT_DEST = "/public/"
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
	SetDebugMode(activate bool)
	SetVerboseMode(activate bool)

	Generate(intput string, output string)
}

// Struct are not exported
type Remap struct {
	url string
	dir string
}

// Define struct to configure and run Prince command line
type Prince struct {
	exePath string

	styleSheets     []string
	scripts         []string
	fileAttachments []string
	remaps          []Remap

	inputType  string
	javascript bool
	insecure   bool

	licenseFile string
	licenseKey  string

	logFile string
	debug   bool
	verbose bool
}

func NewWrapper() Wrapper {

	w := new(Prince)
	w.exePath = PRINCE_BIN

	w.styleSheets = make([]string, 0, 50)
	w.scripts = make([]string, 0, 50)
	w.fileAttachments = make([]string, 0, 50)
	w.remaps = make([]Remap, 0, 50)

	w.SetHTML(true)

	w.debug = true
	w.verbose = true
	return w
}

func (w *Prince) Generate(input string, output string) {

	_, err := exec.LookPath(w.exePath)
	if err != nil {
		// return an error
		return

	}

	outputFile := filepath.Join(OUTPUT_DEST, output)

	args := w.GetCommandLineArgs(outputFile)
	args = append(args, input)

	_, err = exec.Command(w.exePath, args...).Output()

	if nil != err {
		fmt.Println(err.Error())
	}
}

func (w *Prince) GetCommandLineArgs(outputFile string) []string {

	//args := make([]string, 0, cap(w.GetCommandLineArgsCount()))
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

	if "" != w.logFile {
		args = append(args, "--log="+strconv.Quote(w.logFile))
	}
	if true == w.debug {
		args = append(args, "--debug")
	}
	if true != w.verbose {
		args = append(args, "--verbose")
	}

	args = append(args, "--structured-log=normal", "-o", outputFile)

	return args
}

func (w *Prince) GetCommandLineArgsCount() int {

	argsLen := 0
	argsLen += len(w.styleSheets) * 2
	argsLen += len(w.scripts) * 2
	argsLen += len(w.fileAttachments)
	argsLen += len(w.remaps)

	if "auto" != w.inputType {
		argsLen++
	}
	if true == w.javascript {
		argsLen++
	}
	if true == w.insecure {
		argsLen++
	}

	// License
	if "" != w.licenseKey {
		argsLen++
	}
	if "" != w.licenseFile {
		argsLen++
	}

	if "" != w.logFile {
		argsLen++
	}
	if true == w.debug {
		argsLen++
	}
	if true != w.verbose {
		argsLen++
	}

	return int(argsLen)
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

// Define file input type
func (w *Prince) SetHTML(isHTML bool) {
	if isHTML == true {
		w.inputType = "html"
	} else {
		w.inputType = "xml"
	}
}

// Specify a file that Prince should use to log error/warning messages.
// logFile: The filename that Prince should use to log error/warning
// messages, or '' to disable logging.
func (w *Prince) SetLogFile(logPath string) {
	w.logFile = logPath
}

func (w *Prince) SetDebugMode(activate bool) {
	w.debug = activate
}

func (w *Prince) SetVerboseMode(activate bool) {
	w.verbose = activate
}
