package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {

	allowedOrigins := map[string]bool{
		"http://localhost:4321":    true,
		"http://192.168.3.16:1000": true,
	}

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowOriginsFunc: func(origin string) bool {
			return allowedOrigins[origin]
		},
	}))

	api := app.Group("/api")

	v1 := api.Group("/v1")

	// PEDIDOS
	v1.Get("/pedidos", ObtenerPedidos)
	v1.Get("/pedidos/:id", ObtenerPedido)
	v1.Get("/pedidos/por-comprador/:id", ObtenerPedidosPorComprador)
	v1.Put("/pedidos/asignar", AsignarPedido)
	v1.Patch("/pedidos/desasignar/:id", DesasignaPedido)
	v1.Delete("/pedidos/:id", EliminarPedido)

	// COMPRADOR
	v1.Get("/comprador", OptenerCompradores)
	v1.Get("/comprador/:id", OptenerComprador)
	v1.Post("/comprador", CrearComprador)
	v1.Put("/comprador", ActualizarComprador)
	v1.Delete("/comprador/:id", EliminarComprador)
}
