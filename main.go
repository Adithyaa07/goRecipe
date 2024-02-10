package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Recipe struct {
	Title        string `json:"title"`
	Ingredients  string `json:"ingredients"`
	Servings     string `json:"servings"`
	Instructions string `json:"instructions"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("No Key Found. Use Recipe API key")
	}
	apiKey := os.Getenv("RECIPE_API_KEY")
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go QUERY")
		os.Exit(1)
	}

	query := os.Args[1]

	url := fmt.Sprintf("https://api.api-ninjas.com/v1/recipe?query=%s", query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("X-Api-Key", apiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Error: %d %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	var recipes []Recipe
	err = json.Unmarshal(body, &recipes)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}

	if len(recipes) == 0 {
		fmt.Println("No recipes found for the given query:(")
	} else {
		fmt.Printf("Found %d recipes for \"%s\":\n", len(recipes), query)
		fmt.Println("----------------------------------")
		for _, recipe := range recipes {
			fmt.Printf("Title: %s\n\n", recipe.Title)
			fmt.Printf("Ingredients: %s\n\n", recipe.Ingredients)
			fmt.Printf("Servings: %s\n\n", recipe.Servings)
			fmt.Printf("Instructions: %s\n\n", recipe.Instructions)
			fmt.Println("===============================================================================")
		}
	}
}
