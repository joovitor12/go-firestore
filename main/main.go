package main

import (
	"context"
	"log"
	"main/controller"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

type Character struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

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

	// rotas marvel
	controller.MarvelController(appFiber, client)

	// definindo porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// inicializando server
	log.Fatal(appFiber.Listen(":" + port))

}
