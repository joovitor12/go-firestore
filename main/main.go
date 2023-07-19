package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	// Defina uma rota para ler os dados do Firestore
	appFiber.Get("/", func(c *fiber.Ctx) error {
		doc, err := client.Collection("marvel-characters").Doc("character").Get(context.Background())
		if err != nil {
			return c.Status(500).SendString("Erro ao ler dados do Firestore")
		}
		return c.JSON(doc.Data())
	})

	// Defina uma rota para escrever dados no Firestore
	appFiber.Post("/marvel/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")

		// Parseie o corpo da requisição para a struct Character
		character := findMarvelCharacters(name, c)
		fmt.Print(character)
		if err := c.BodyParser(character); err != nil {
			return c.Status(400).SendString("Dados inválidos")
		}

		// Chame a função para buscar o personagem na API da Marvel
		err := findMarvelCharacters(name, c)
		if err != nil {
			return c.Status(500).SendString("Erro ao buscar personagem na API da Marvel")
		}

		// Escreva os dados do Character no Firestore com o ID como chave do documento
		_, err = client.Collection("characters").Doc("content").Set(context.Background(), character)
		if err != nil {
			return c.Status(500).SendString("Erro ao escrever dados no Firestore")
		}

		return c.SendString("Personagem armazenado no Firestore com sucesso")
	})
	// Defina a porta na qual o servidor irá ouvir
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Inicie o servidor Fiber
	log.Fatal(appFiber.Listen(":" + port))

}

func findMarvelCharacters(name string, c *fiber.Ctx) error {
	response, err := http.Get("http://gateway.marvel.com/v1/public/characters?name=" + name + "&ts=1&apikey=5ee728ad5618f7807e45d5f757e08697&hash=711c3b1a81ec2711ec846e2ab60e91c8")

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error at API: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error at body: %s", err)
	}

	var result struct {
		Data struct {
			Results []Character `json:"results"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Error at Unmarshal: %s", err)
	}

	if len(result.Data.Results) == 0 {
		return c.Status(http.StatusNotFound).SendString("Personagem não encontrado")
	}

	character := result.Data.Results[0]

	fmt.Printf("ID: %d\n", character.ID)
	fmt.Printf("Name: %s\n", character.Name)
	fmt.Printf("Description: %s\n", character.Description)

	//esse retorno ta errado, tem que retornar o character e nao o json dele
	return c.JSON(character)
}
