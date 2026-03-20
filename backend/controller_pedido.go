package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func ObtenerPedidos(c *fiber.Ctx) error {

	db := GetDB()

	rows, err := db.Query(`
		SELECT 
			p.id_pedido, 
			p.nombre AS pedido, 
			p.precio, 
			p.imagen, 
			p.variante, 
			p.cantidad, 
			COALESCE(c.nombre, '') AS comprador,
			COALESCE(c.id_comprador, 0) AS id_comprador
		FROM pedido AS p
			LEFT JOIN comprador AS c ON p.id_comprador = c.id_comprador
		ORDER BY p.id_pedido
	`)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error consultando pedidos",
		})
	}

	defer rows.Close()

	pedidos := []PedidosDTO{}

	for rows.Next() {

		var p PedidosDTO

		if err := rows.Scan(
			&p.Id_Pedido,
			&p.Nombre,
			&p.Precio,
			&p.Imagen,
			&p.Variante,
			&p.Cantidad,
			&p.Comprador,
			&p.Id_Comprador,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error leyendo pedidos",
			})
		}

		pedidos = append(pedidos, p)

	}

	return c.JSON(pedidos)
}

func ObtenerPedido(c *fiber.Ctx) error {
	db := GetDB()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID requerido",
		})
	}

	var pedido PedidosDTO

	err := db.QueryRow(`
		SELECT 
			p.id_pedido, 
			p.nombre AS pedido, 
			p.precio, 
			p.imagen, 
			p.variante, 
			p.cantidad, 
			COALESCE(c.nombre, '') AS comprador,
			COALESCE(c.id_comprador, 0) AS id_comprador
		FROM pedido p
		LEFT JOIN comprador c ON p.id_comprador = c.id_comprador
		WHERE p.id_pedido = ?
		LIMIT 1
	`, id).Scan(
		&pedido.Id_Pedido,
		&pedido.Nombre,
		&pedido.Precio,
		&pedido.Imagen,
		&pedido.Variante,
		&pedido.Cantidad,
		&pedido.Comprador,
		&pedido.Id_Comprador,
	)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pedido no encontrado",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error obteniendo pedido",
		})
	}

	return c.JSON(pedido)
}

func ObtenerPedidosPorComprador(c *fiber.Ctx) error {
	db := GetDB()

	idComprador := c.Params("id")
	if idComprador == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID del comprador requerido",
		})
	}

	rows, err := db.Query(`
		SELECT 
			p.id_pedido, 
			p.nombre AS pedido, 
			p.precio, 
			p.imagen, 
			p.variante, 
			p.cantidad, 
			COALESCE(c.nombre, '') AS comprador,
			COALESCE(c.id_comprador, 0) AS id_comprador
		FROM pedido p
		INNER JOIN comprador c ON p.id_comprador = c.id_comprador
		WHERE p.id_comprador = ?
		ORDER BY p.id_pedido
	`, idComprador)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error consultando pedidos por comprador",
		})
	}
	defer rows.Close()

	pedidos := []PedidosDTO{}

	for rows.Next() {
		var p PedidosDTO

		if err := rows.Scan(
			&p.Id_Pedido,
			&p.Nombre,
			&p.Precio,
			&p.Imagen,
			&p.Variante,
			&p.Cantidad,
			&p.Comprador,
			&p.Id_Comprador,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error leyendo pedidos",
			})
		}

		pedidos = append(pedidos, p)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error iterando pedidos",
		})
	}

	return c.JSON(pedidos)
}

func AsignarPedido(c *fiber.Ctx) error {

	db := GetDB()

	var req AsignarPedidoRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	if req.Id_Pedido <= 0 || req.Id_Comprador <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id_pedido e id_comprador son obligatorios",
		})
	}

	result, err := db.Exec(`
		UPDATE pedido
		SET id_comprador = ?
		WHERE id_pedido = ?
	`, req.Id_Comprador, req.Id_Pedido)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error asignando pedido",
		})
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pedido no encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pedido asignado correctamente",
	})

}

func DesasignaPedido(c *fiber.Ctx) error {
	db := GetDB()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID del pedido requerido",
		})
	}

	result, err := db.Exec(`
		UPDATE pedido
		SET id_comprador = NULL
		WHERE id_pedido = ?
	`, id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error desasignando pedido",
		})
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo verificar la desasignación",
		})
	}

	if rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pedido no encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pedido desasignado correctamente",
	})
}

func EliminarPedido(c *fiber.Ctx) error {
	db := GetDB()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID del pedido requerido",
		})
	}

	result, err := db.Exec(`
		DELETE FROM pedido
		WHERE id_pedido = ?
	`, id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error eliminando pedido",
		})
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo verificar la eliminación",
		})
	}

	if rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pedido no encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pedido eliminado correctamente",
	})
}
