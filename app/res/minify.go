package res

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func Minify() error {
	directory := "web/public"
	minJsFile := MinJsFile
	minJsFullPath := directory + "/" + minJsFile
	fileList := getJsFilesToMinify()

	var possibleCommands = []string{
		"/usr/local/bin/yuicompressor",
		"yui-compressor",
	}
	var yuiCommand string
	for _, possibleCommand := range possibleCommands {
		if isCommandAvailable(possibleCommand) {
			yuiCommand = possibleCommand
			break
		}
	}
	if yuiCommand == "" {
		return jerr.New("yui command not found")
	}

	var wg sync.WaitGroup
	var errors []error
	var uncompressedSize int64
	compressedData := make(map[string]string)
	headFormat := "%-40s %12s %12s %12s\n"
	lineFormat := "%-40s %12d %12d %11.2f%%\n"

	fmt.Println("Building...")
	fmt.Printf(headFormat, "FILE", "LEN", "COMPRESSED", "RATIO")

	wg.Add(len(fileList))
	for _, filename := range fileList {
		go func(filename string) {
			fullFilename := directory + "/" + filename
			defer wg.Done()
			file, err := os.Open(fullFilename)
			if err != nil {
				errors = append(errors, jerr.Get("error opening file", err))
				return
			}
			fi, err := file.Stat()
			if err != nil {
				errors = append(errors, jerr.Get("error getting file stat", err))
				return
			}
			uncompressedSize += fi.Size()

			out, err := exec.Command("sh", "-c", fmt.Sprintf("%s %s", yuiCommand, fullFilename)).CombinedOutput()
			if err != nil {
				errors = append(errors, jerr.Get("error running yui-compressor", err))
				return
			}
			compressedData[filename] = string(out)
			fmt.Printf(lineFormat, filename, fi.Size(), len(out), (1-(float32(len(out)))/float32(fi.Size()))*100)
		}(filename)
	}
	wg.Wait()

	if len(errors) > 0 {
		return jerr.Get("error building js", jerr.Combine(errors...))
	}

	fmt.Println()

	var compressedDataSorted []string
	for _, filename := range fileList {
		compressedDataSorted = append(compressedDataSorted, compressedData[filename])
	}

	jsMin := []byte(strings.Join(compressedDataSorted, "\n"))
	ioutil.WriteFile(minJsFullPath, jsMin, 0644)
	msg := "Wrote: %s, Len: %d, Uncompressed: %d, Compression ratio: %0.2f%%\n"
	fmt.Printf(msg, minJsFullPath, len(jsMin), uncompressedSize, (1-(float32(len(jsMin))/float32(uncompressedSize)))*100)
	return nil
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("command -v '%s'", name))
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
