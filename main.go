package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func main() {
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

	app.Listen(":8000")
}

type Person struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
