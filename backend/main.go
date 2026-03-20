package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	InitDB()

	defer GetDB().Close()

	SetupRoutes(app)

	_ = app.Listen(fmt.Sprintf(":%s", "5000"))

}
