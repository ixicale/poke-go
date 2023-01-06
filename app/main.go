package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const DNS = "https://pokeapi.co"

func main() {
	router := gin.Default()
	router.GET("/pokemons", func(c *gin.Context) {
		pokemons := GetAllPokemons()

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Pokemons List",
			"data":    pokemons,
		})
	})

	router.GET("/pokemons/:name", func(c *gin.Context) {
		name := c.Param("name")
		data := getPokemonByName(name)

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Pokemon",
			"data":    data,
		})
	})

	router.Run()
}

func make_request(sufix string) (map[string]interface{}, error) {
	resp, err := http.Get(DNS + sufix)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}

func GetAllPokemons() []string {
	sufix := "/api/v2/pokemon?limit=1000"
	response, _ := make_request(sufix)

	var pokemons []string
	// iterate over the response.results
	for _, p := range response["results"].([]interface{}) {
		pokemons = append(pokemons, p.(map[string]interface{})["name"].(string))
	}

	return pokemons
}

func getPokemonByName(name string) map[string]interface{} {
	sufix := fmt.Sprintf("/api/v2/pokemon/%s", name)
	response, _ := make_request(sufix)

	return response
}
