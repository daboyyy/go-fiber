package main

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"github.com/gofiber/fiber/v2"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	jwtware "github.com/gofiber/jwt/v2"
)

var db *sqlx.DB

const jwtSecret = "infinitas"

func main() {
	var err error
	db, err = sqlx.Open("mysql", "root:P@ssw0rd@tcp(13.76.163.73:3306)/techcoach")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use("/hello", jwtware.New(jwtware.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(jwtSecret),
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return fiber.ErrUnauthorized
		},
	}))

	app.Post("/signup", Signup)
	app.Post("/login", Login)
	app.Get("/hello", Hello)

	app.Listen(":8000")
}

func Signup(c *fiber.Ctx) error {
	request := SignupRequest{}
	err := c.BodyParser(&request)
	if err != nil{
		return err
	}

	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	query := "insert user (username, password) values (?, ?)"
	result, err := db.Exec(query, request.Username, string(password))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	user := User{
		Id: int(id),
		Username: request.Username,
		Password: string(password),
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func Login(c *fiber.Ctx) error {
	request := LoginRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	user := User{}
	query := "select id, username, password from user where username=?"
	err = db.Get(&user, query, request.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Incorrect username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Incorrect username or password")
	}

	cliams := jwt.StandardClaims{
		Issuer:    strconv.Itoa(user.Id),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	token, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"jwtToken": token,
	})
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

type User struct {
	Id int `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
