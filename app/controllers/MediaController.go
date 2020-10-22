package controllers

import (
	"fmt"
	"golang-api/app/helpers"
	"golang-api/app/models"
	"golang-api/app/structs"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaController struct{}

func (controller MediaController) Store(c *gin.Context) {

	var failedValidations = map[string]interface{}{}

	file, _ := c.FormFile("media")
	if file == nil {
		failedValidations["media"] = []string{"The media field is required."}

		c.JSON(422, helpers.Validator{
			Message: "The given data was invalid",
			Errors:  failedValidations,
		})
		return
	}

	fullFilename := uuid.New().String() + filepath.Ext(file.Filename)

	media := structs.Media{
		Filename: fullFilename,
		Type:     file.Header["Content-Type"][0],
		Location: "https://" + os.Getenv("AWS_BUCKET_NAME") + ".s3.amazonaws.com/" + fullFilename,
	}

	fileBuffer, err := file.Open()
	if err != nil {
		fmt.Println("Error opening buffer of uploaded file", err)
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_DEFAULT_REGION"))})
	if err != nil {
		log.Fatal(err)
	}

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(media.Filename),
		ACL:    aws.String("public-read"),
		Body:   fileBuffer,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	media, err = models.StoreMedia(media)
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	c.JSON(200, media.Response())
}
