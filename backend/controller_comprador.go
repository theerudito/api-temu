package main

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func OptenerCompradores(c *fiber.Ctx) error {
	db := GetDB()

	rows, err := db.Query(`
		SELECT id_comprador, nombre
		FROM comprador
		ORDER BY id_comprador
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error consultando compradores",
		})
	}
	defer rows.Close()

	compradores := []Comprador{}

	for rows.Next() {
		var comp Comprador
		if err := rows.Scan(
			&comp.Id_Comprador,
			&comp.Nombre,
		); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error leyendo compradores",
			})
		}
		compradores = append(compradores, comp)
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error iterando compradores",
		})
	}

	return c.JSON(compradores)
}

func OptenerComprador(c *fiber.Ctx) error {
	db := GetDB()

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "ID requerido",
		})
	}

	var comp Comprador

	err := db.QueryRow(`
		SELECT id_comprador, nombre
		FROM comprador
		WHERE id_comprador = ?
	`, id).Scan(
		&comp.Id_Comprador,
		&comp.Nombre,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error": "Comprador no encontrado",
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error obteniendo comprador",
		})
	}

	return c.JSON(comp)
}

func CrearComprador(c *fiber.Ctx) error {
	db := GetDB()

	var comp Comprador

	if err := c.BodyParser(&comp); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	if comp.Nombre == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "El nombre del comprador es obligatorio",
		})
	}

	result, err := db.Exec(`
		INSERT INTO comprador (nombre)
		VALUES (?)
	`, strings.ToUpper(comp.Nombre))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error creando comprador",
		})
	}

	id, _ := result.LastInsertId()
	comp.Id_Comprador = int(id)

	return c.Status(201).JSON(comp)
}

func ActualizarComprador(c *fiber.Ctx) error {

	db := GetDB()

	var comp Comprador
	if err := c.BodyParser(&comp); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	if comp.Nombre == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "El nombre del comprador es obligatorio",
		})
	}

	result, err := db.Exec(`
		UPDATE comprador
		SET nombre = ?
		WHERE id_comprador = ?
	`, strings.ToUpper(comp.Nombre), comp.Id_Comprador)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error actualizando comprador",
		})
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Comprador no encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comprador actualizado correctamente",
	})
}

func EliminarComprador(c *fiber.Ctx) error {
	db := GetDB()

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "ID requerido",
		})
	}

	result, err := db.Exec(`
		DELETE FROM comprador
		WHERE id_comprador = ?
	`, id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error eliminando comprador",
		})
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Comprador no encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comprador eliminado correctamente",
	})
}
