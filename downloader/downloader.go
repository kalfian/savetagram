package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints the progress of a file write
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)

	fmt.Printf("\rDownloading... %s complete", humanize.IBytes(wc.Total))
}

func DownloadFile(url string) (string, error) {

	fileName := ""
	fileTemp := strconv.Itoa(int(time.Now().Unix()))
	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	out, err := os.Create(fileTemp + ".tmp")
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	fileName = path.Base(resp.Request.URL.Path)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create our bytes counter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return "", err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Println()

	// Rename the tmp file back to the original file
	err = os.Rename(fileTemp+".tmp", fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
