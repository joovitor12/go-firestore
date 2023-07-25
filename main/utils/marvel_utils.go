package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Character struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetMarvelCharacter(name string, c *fiber.Ctx) (*Character, error) {
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
		return nil, fmt.Errorf("Personagem não encontrado")
	}

	character := result.Data.Results[0]

	fmt.Printf("ID: %d\n", character.ID)
	fmt.Printf("Name: %s\n", character.Name)
	fmt.Printf("Description: %s\n", character.Description)

	return &character, nil
}

func GetMarvelCharacterFromDB(c *fiber.Ctx, name string) (*Character, error) {
	client, ok := c.Locals("firebase").(*firestore.Client)
	if !ok {
		log.Println("Firebase client not found in context locals.")
		return nil, fmt.Errorf("Internal server error")
	}

	// Usando o próprio contexto do Fiber para as operações com Firestore
	ctx := c.Context()

	// Acesso à coleção "marvel-characters" e ao documento com o nome do personagem
	doc, err := client.Collection("marvel-characters").Doc(name).Get(ctx)
	if err != nil {
		// Verifica se o documento não foi encontrado no Firestore
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("Character not found")
		}
		log.Printf("Error in getting character: %v", err)
		return nil, fmt.Errorf("Internal server error")
	}

	// Cria uma variável para armazenar os dados do documento
	var character Character

	// Converte os dados do documento para a estrutura Character
	if err := doc.DataTo(&character); err != nil {
		log.Printf("Error in getting character from Firestore: %v - %v", err, character)
		return nil, fmt.Errorf("Internal server error")
	}

	// Imprime algumas informações do personagem no log para fins de depuração
	log.Printf("ID: %d\n", character.ID)
	log.Printf("Name: %s\n", character.Name)
	log.Printf("Description: %s\n", character.Description)

	return &character, nil
}
