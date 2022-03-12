package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// GET
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("GET: Hello World")
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

	app.Listen(":8000")
}

type Person struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
