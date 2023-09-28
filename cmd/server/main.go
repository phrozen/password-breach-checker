package main

import (
	_ "embed"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/phrozen/password-breach-checker/pkg/database"
)

//go:embed index.html
var index []byte

type CheckRequest struct {
	Hash string `json:"hash" xml:"hash" form:"hash"`
}

func main() {

	input := flag.String("db", "", "Input filename of passwords ordered by hash")
	port := flag.String("port", "3000", "Port to run server")
	logs := flag.Bool("log", true, "Enable/disable logging middleware")
	flag.Parse()

	db, err := database.New(*input)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	app := fiber.New()

	if *logs {
		app.Use(logger.New())
	}

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.Send(index)
	})

	app.Post("/check", func(c *fiber.Ctx) error {
		var req CheckRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		hash, err := hex.DecodeString(req.Hash)
		if err != nil {
			return err
		}
		return c.SendString(fmt.Sprintf("%d", db.Search(hash)))
	})

	app.Listen(":" + *port)
}
