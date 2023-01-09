package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"

	"github.com/gin-gonic/gin"
)

const DNS = "https://pokeapi.co"

// var ginLambda *ginadapter.GinLambda

func main() {
	router := gin.New()
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

		reqStatus := http.StatusOK
		data := getPokemonByName(name)
		var message string
		// check if name is number
		if _, err := strconv.Atoi(name); err == nil {
			reqStatus = http.StatusBadRequest
			message = "Name must be a string"
			data = nil
		} else if data == nil && reqStatus == http.StatusOK {
			reqStatus = http.StatusNotFound
			message = fmt.Sprintf("Pokemon '%s' not found", name)
		} else {
			message = fmt.Sprintf("Pokemon %s", name)
		}

		c.JSON(reqStatus, gin.H{
			"status":  reqStatus,
			"message": message,
			"data":    data,
		})
	})

	router.Run()
	// ginLambda = ginadapter.New(router)
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
