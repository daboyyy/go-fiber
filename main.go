package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"time"
)

func main() {

}

func Fiber() {
	app := fiber.New()

	// Middleware (every endpoing)
	app.Use(func(c *fiber.Ctx) error {
		fmt.Println("before")
		err := c.Next()
		fmt.Println("after")
		return err
	})

	// Middleware (specific endpoing)
	app.Use("/hello", func(c *fiber.Ctx) error {
		c.Locals("name", "bond")
		fmt.Println("before")
		c.Next()
		fmt.Println("after")
		return nil
	})

	// Middleware Requestid
	app.Use(requestid.New())

	// Cors
	app.Use(cors.New())

	// Logger
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))

	// GET
	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Locals("name")
		return c.SendString(fmt.Sprintf("GET: Hello %v", name))
	})

	// POST
	app.Post("/hello", func(c *fiber.Ctx) error {
		return c.SendString("POST: Hello World")
	})

	// Parameters
	app.Get("/hello/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString("name: " + name)
	})

	// Parameters Optional
	app.Get("/hello/:name/:surname?", func(c *fiber.Ctx) error {
		name := c.Params("name")
		surname := c.Params("surname") // optional param
		return c.SendString("name: " + name + ", surname: " + surname)
	})

	// ParamsInt
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.ErrBadRequest
		}
		return c.SendString(fmt.Sprintf("ID = %v", id))
	})

	// Query
	// localhost:8000/query?name=bond
	app.Get("/query", func(c *fiber.Ctx) error {
		name := c.Query("name")
		surname := c.Query("surname")
		return c.SendString("name: " + name + ", surname: " + surname)
	})

	// QueryParser
	app.Get("/query2", func(c *fiber.Ctx) error {
		person := Person{}
		c.QueryParser(&person)
		return c.JSON(person)
	})

	// Wildcards
	app.Get("/wildcards/*", func(c *fiber.Ctx) error {
		wildcard := c.Params("*")
		return c.SendString(wildcard)
	})

	// Static File
	app.Static("/", "./wwwroot", fiber.Static{
		Index: "index.html",
		CacheDuration: time.Second * 10,
	})

	// NewError
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "content not found")
	})

	// Group
	v1 := app.Group("/v1")
	v1.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello v1")
	})

	v2 := app.Group("/v2", func(c *fiber.Ctx) error {
		c.Set("version", "v2")
		return c.Next()
	})
	v2.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello v2")
	})

	// Mount
	userApp := fiber.New()
	userApp.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login")
	})
	app.Mount("/user", userApp)

	// Server
	app.Server().MaxConnsPerIP = 1
	app.Get("/server", func(c *fiber.Ctx) error {
		time.Sleep(time.Second * 30)
		return c.SendString("Server")
	})

	// env
	app.Get("/env", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"BaseURL": c.BaseURL(),
			"Hostname": c.Hostname(),
			"IP": c.IP(),
			"IPs": c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path": c.Path(),
			"Protocol": c.Protocol(),
			"SubDomains": c.Subdomains(),
		})
	})

	// Body
	app.Post("/body", func(c *fiber.Ctx) error {
		fmt.Printf("IsJson: %v \n", c.Is("json"))
		fmt.Println(string(c.Body()))

		person := Person{}
		err := c.BodyParser(&person)
		if err != nil {
			return err
		}

		fmt.Println(person)
		return nil
	})

	app.Post("/body2", func(c *fiber.Ctx) error {
		fmt.Printf("IsJson: %v \n", c.Is("json"))
		fmt.Println(string(c.Body()))

		data := map[string]interface{}{}
		err := c.BodyParser(&data)
		if err != nil {
			return err
		}

		fmt.Println(data)
		return nil
	})

	app.Listen(":8000")
}

type Person struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
