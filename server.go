package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Paste struct {
	Content string `form:"content" binding:"required"`
}

type UploadFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func paste(c *gin.Context) {
	// ipfs gateway url
	gatewayURL := "http://localhost:8080"

	// json data
	type ContentRequestBody struct {
		Content string `json:"content"`
	}

	var requestBody ContentRequestBody

	c.BindJSON(&requestBody)

	// Write content to a file
	filename := time.Now().Format(time.UnixDate) + ".txt"
	filepath := path.Join("files", "text", filename)

	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(requestBody.Content)
	f.Close()

	// Add this file to IPFS
	cmd := exec.Command("ipfs", "add", filepath)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// Create a share URL and return
	words := strings.Split(string(output[:]), " ")
	hash := words[1]
	url := gatewayURL + "/ipfs/" + hash
	response := map[string]interface{}{"url": url}

	c.JSONP(http.StatusOK, response)
}

func upload(c *gin.Context) {
	// ipfs gateway url
	gatewayURL := "http://localhost:8080"

	infile, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	// Create a directory to put this file in
	hasher := md5.New()
	io.WriteString(hasher, header.Filename)
	dirname := hex.EncodeToString(hasher.Sum(nil))[:5]
	if err := os.MkdirAll(path.Join("files", "other", dirname), 0755); err != nil {
		log.Fatal(err)
	}

	outfilename := header.Filename
	outfilepath := path.Join("files", "other", dirname, outfilename)

	outfile, err := os.Create(outfilepath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Piping upload into disk file.")
	io.Copy(outfile, infile)
	outfile.Close()

	// Add to IPFS
	cmd := exec.Command("ipfs", "add", "-r", path.Join("files", "other", dirname))
	output, err := cmd.Output()
	log.Println(string(output))
	if err != nil {
		log.Fatal(err)
	}

	// Create a File IPFS and return
	lines := strings.Split(string(output[:]), "\n")
	words := strings.Split(lines[len(lines)-2], " ")
	hash := words[1]
	url := gatewayURL + "/ipfs/" + hash + "/" + outfilename
	response := map[string]interface{}{"url": url}

	c.JSONP(http.StatusOK, response)
}
