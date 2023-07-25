package controller

import (
	"context"
	"fmt"
	"main/utils"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
)

func MarvelController(app *fiber.App, client *firestore.Client) {

	app.Use(setupFirestore(client))

	app.Get("/", func(c *fiber.Ctx) error {
		doc, err := client.Collection("marvel-characters").Doc("character").Get(c.Context())
		if err != nil {
			return c.Status(500).SendString("Erro ao ler dados do Firestore")
		}
		return c.JSON(doc.Data())
	})

	app.Post("/marvel/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")

		// firestoreObj, err := utils.GetMarvelCharacterFromDB(c, name)

		// if err != nil {
		// 	return c.Status(500).SendString("Erro ao buscar personagem no Firestore")
		// }

		// if firestoreObj != nil {
		// 	return c.JSON(firestoreObj)
		// }

		// buscando o personagem na API da Marvel
		character, err := utils.GetMarvelCharacter(name, c)
		if err != nil {
			return c.Status(500).SendString("Erro ao buscar personagem na API da Marvel")
		}

		// escrevendo dados no Firestore
		marvelCollectionRef := client.Collection("marvel-characters")
		_, err = marvelCollectionRef.Doc(fmt.Sprintf("%d", character.ID)).Set(context.Background(), character)
		if err != nil {
			return c.Status(500).SendString("Erro ao escrever dados no Firestore")
		}

		return c.SendString("Personagem armazenado no Firestore com sucesso")
	})
}

func setupFirestore(firestoreClient *firestore.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("firebase", firestoreClient)
		return c.Next()
	}
}
