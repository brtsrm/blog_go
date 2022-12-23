package controller

import (
	"blog/database"
	"blog/models"
	"blog/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile("[a-z0-9.%+-]+@[a-z09.%+-]+.[a-z]")
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}
	// Check if password is less than 6 characters
	if len(data["password"].(string)) <= 6 {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "PAssword must be greater than 6 character",
		})
	}
	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid Email Address",
		})
	}
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id != 0 {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Email already exits",
		})
	}

	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		Phone:     data["phone"].(string),
	}
	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)

	}
	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Account created successfully",
	})
}
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}
	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Email address doesn't exit, kindly create an account",
		})
	}
	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}
	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return nil
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "You have successfully login",
		"user":    user,
	})

}

type Claims struct {
	jwt.StandardClaims
}
