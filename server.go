package main

import (
	"log"
	"mime/multipart"
	"net/http"

	"github.com/IPFS-NEXIVIL/orbit-db-gateway/database"
	"github.com/IPFS-NEXIVIL/orbit-db-gateway/models"
	"github.com/gin-gonic/gin"
)

type DBInfo struct {
	DB *database.Database
}

type Paste struct {
	Content string `form:"content" binding:"required"`
}

type UploadFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (database DBInfo) saveAndGetDBData(data string) models.Data {
	// our orbit db
	db := database.DB

	// new data create
	newData := models.NewData()
	log.Println(newData)

	newData.Content = data

	// insert `content` data to orbit db
	db.SubmitData(newData)

	nexivilData, err := db.GetDataByID(newData.ID)
	if err != nil {
		log.Fatal(err)
	}

	return nexivilData
}

func (database DBInfo) paste(c *gin.Context) {
	// json data
	type ContentRequestBody struct {
		Content string `json:"content"`
	}

	var requestBody ContentRequestBody

	newData := models.NewData()
	log.Println(newData)

	c.BindJSON(&requestBody)

	data := database.saveAndGetDBData(requestBody.Content)

	c.String(http.StatusOK, "data %s save to orbit db success", data.Content)
}
