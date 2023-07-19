package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

func main() {
	// Inicialize o Firebase com suas credenciais
	opt := option.WithCredentialsFile("serviceAccountKeyPDI.json")
	config := &firebase.Config{ProjectID: "pdi-go"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("Erro ao inicializar o app Firebase: %v\n", err)
	}

	// Inicialize o Firestore
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Erro ao inicializar o Firestore: %v\n", err)
	}
	defer client.Close()

	// Crie um novo aplicativo Fiber
	appFiber := fiber.New()

	// Defina uma rota para ler os dados do Firestore
	appFiber.Get("/", func(c *fiber.Ctx) error {
		doc, err := client.Collection("messages").Doc("message").Get(context.Background())
		if err != nil {
			return c.Status(500).SendString("Erro ao ler dados do Firestore")
		}
		return c.JSON(doc.Data())
	})

	// Defina uma rota para escrever dados no Firestore
	appFiber.Post("/", func(c *fiber.Ctx) error {
		data := new(struct {
			Message string `json:"message"`
		})
		if err := c.BodyParser(data); err != nil {
			return c.Status(400).SendString("Dados inválidos")
		}

		_, err := client.Collection("messages").Doc("text").Set(context.Background(), data)
		if err != nil {
			return c.Status(500).SendString("Erro ao escrever dados no Firestore")
		}
		return c.SendString("Dados escritos no Firestore com sucesso")
	})

	// Defina a porta na qual o servidor irá ouvir
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Inicie o servidor Fiber
	log.Fatal(appFiber.Listen(":" + port))
}
