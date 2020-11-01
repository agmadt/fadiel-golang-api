package controllers

import (
	"database/sql"
	"fmt"
	"golang-api/app/helpers"
	"golang-api/app/models"
	"golang-api/app/structs"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (controller AuthController) Login(c *gin.Context) {

	var loginRequest structs.LoginRequest
	var failedValidations map[string]interface{}

	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		c.JSON(400, gin.H{"message": "Something wrong with the request"})
		fmt.Println("Login bind json error", err)
		return
	}

	validate = validator.New()
	err = validate.Struct(loginRequest)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errFailedField := helpers.ValidatorRemoveNamespace(strcase.ToSnake(err.Namespace()))
			failedValidations[errFailedField] = []string{helpers.ValidatorMessage(errFailedField, err.ActualTag(), err.Param())}
		}

		c.JSON(422, helpers.Validator{
			Message: "The given data was invalid",
			Errors:  failedValidations,
		})
		return
	}

	user := structs.User{Email: loginRequest.Email}

	user, err = models.FindUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Username or password is incorrect"})
		} else {
			c.JSON(500, gin.H{"message": "Server error"})
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		c.JSON(404, gin.H{"message": "Username or password is incorrect"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 10).Unix(),
		"iat": time.Now().Unix(),
		"sub": user.ID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		fmt.Println("Error when creating jwt token", err)
		return
	}

	loginResponse := structs.LoginResponse{
		Token: tokenString,
		Data: structs.LoginResponseData{
			ID:    user.ID,
			Email: loginRequest.Email,
			Name:  user.Name,
		},
	}

	c.JSON(200, loginResponse)
}
