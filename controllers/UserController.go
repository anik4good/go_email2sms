package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/anik4good/fiber_boilerplate/util"
	"github.com/bxcodec/faker/v3"

	//"strings"
	"time"

	Configuration "github.com/anik4good/fiber_boilerplate/config"
	"github.com/anik4good/fiber_boilerplate/models"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid"
)

var database *sql.DB

//Hello
func Hello(c *fiber.Ctx) error {
	return c.SendString("fiber")
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	Configuration.GormDBConn.Find(&users)
	return c.JSON(users)
}

// get user from database and return
func GetSingleUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(models.User)
	res := Configuration.GormDBConn.First(&user, id)
	if res.Error != nil {
		return c.Status(404).JSON("No Record Found")
	}
	return c.Status(200).JSON(user)
}

func CreateUser(c *fiber.Ctx) error {

	//	data, _ := json.MarshalIndent(app.Stack(), "", "  ")
	//	fmt.Println(string(data))
	requestBody := c.Body()
	var email models.User
	json.Unmarshal(requestBody, &email)
	_, err := database.Exec(`INSERT INTO users(name, email,status) VALUES (?,?,?)`, email.Name, email.Email, email.Status)
	if err != nil {

		//	panic(err)

		fmt.Println("error creating user:", email.Name)
		json.NewEncoder(c).Encode("error creating user:")
		return nil
		//	json.NewEncoder(c).Encode("error creating user:")
	}

	json.NewEncoder(c).Encode("received Email: " + email.Email)
	return nil
}

//AddBook
func AddUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	Configuration.GormDBConn.Create(&user)
	log.Println("User Created successfully")
	return c.Status(200).JSON(user)
}

//AddBook
func AddUserRandom(c *fiber.Ctx) error {
	user := new(models.User)

	//user.ID = 0
	user.Name = faker.Name()
	user.Email = faker.Email()
	user.Status = 0
	Configuration.GormDBConn.Create(&user)
	log.Println("Random User Created successfully")
	return c.Status(200).JSON(user)
}

func SmsApi(c *fiber.Ctx) error {

	data := new(models.Api_body)
	time.Sleep(1 * time.Second)

	if err := c.BodyParser(data); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	//Validation
	if len(data.Apikey) == 0 || len(data.MessageType) == 0 || len(data.Contacts) == 0 || len(data.Message) == 0 || len(data.Senderid) == 0 {
		return c.Status(400).JSON("data cannot be empty")
	}

	fmt.Println(len(data.Contacts))
	//if len(data.Contacts) !=11 {
	//	return c.Status(400).JSON("Mobile No should be 11 digit")
	//}

	valid := util.GetValidPhoneNumberUpdated(data.Contacts)
	if valid == false {
		return c.Status(400).JSON("Invalid Phone Number")
	}

	id, _ := gonanoid.Nanoid(13)
	response := "Your SMS is Submitted. ID: " + data.Senderid + "_" + id
	log.Println(response)

	return c.Status(200).JSON(response)

}
