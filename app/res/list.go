package res

import (
	"fmt"
	"strings"
)

var JsFiles = []string{
	"https://cdnjs.cloudflare.com/ajax/libs/jquery/3.3.1/jquery.min.js",
	"https://cdnjs.cloudflare.com/ajax/libs/jqueryui/1.12.1/jquery-ui.min.js",
	"https://cdnjs.cloudflare.com/ajax/libs/identicon.js/2.3.2/pnglib.min.js",
	"https://cdnjs.cloudflare.com/ajax/libs/identicon.js/2.3.2/identicon.min.js",
	"https://cdnjs.cloudflare.com/ajax/libs/jstimezonedetect/1.0.6/jstz.min.js",
	"https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.7/js/bootstrap.min.js",
	"js/init.js",
	"js/login.js",
	"js/signup.js",
	"js/key.js",
	"js/memo.js",
	"js/topics.js",
	"js/profile.js",
	"js/poll.js",
	"js/vote.js",
	"js/modal.js",
	"js/mini-profile.js",
}

var CssFiles = []string{
	"https://cdnjs.cloudflare.com/ajax/libs/jqueryui/1.12.1/jquery-ui.min.css",
	"https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.7/css/bootstrap.min.css",
	"https://cdnjs.cloudflare.com/ajax/libs/flag-icon-css/0.8.2/css/flag-icon.min.css",
	"style.css",
}

var MinJsFile = "js/min.js"

var appendNumber = 0

func SetAppendNumber(num int) {
	appendNumber = num
}

func GetResCssFiles() []string {
	var fileList []string
	for _, file := range CssFiles {
		if ! strings.HasPrefix(file, "lib/") && ! strings.HasPrefix(file, "http") && appendNumber > 0 {
			file = fmt.Sprintf("%s?ver=%d", file, appendNumber)
		}
		fileList = append(fileList, file)
	}
	return fileList
}

func GetResJsFiles() []string {
	var fileList []string
	for _, file := range JsFiles {
		if ! strings.HasPrefix(file, "lib/") && ! strings.HasPrefix(file, "http") && appendNumber > 0 {
			file = fmt.Sprintf("%s?ver=%d", file, appendNumber)
		}
		fileList = append(fileList, file)
	}
	return fileList
}

func GetMinJsFiles() []string {
	var fileList []string
	for _, file := range JsFiles {
		if strings.HasPrefix(file, "lib/") || strings.HasPrefix(file, "http") {
			fileList = append(fileList, file)
		}
	}
	fileList = append(fileList, fmt.Sprintf("%s?ver=%d", MinJsFile, appendNumber))
	return fileList
}

func getJsFilesToMinify() []string {
	var fileList []string
	for _, file := range JsFiles {
		if ! strings.HasPrefix(file, "lib/") && ! strings.HasPrefix(file, "http") {
			fileList = append(fileList, file)
		}
	}
	return fileList
}
