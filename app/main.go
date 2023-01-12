package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"reflect"

	"github.com/gin-gonic/gin"
)

const DNS = "https://pokeapi.co"

func main() {
	// LOCAL_DEBUG is true when AWS_EXECUTION_ENV is empty
	LOCAL_DEBUG := reflect.DeepEqual(os.Getenv("AWS_EXECUTION_ENV"), "")

	router := gin.New()
	router.GET("/pokemon", func(c *gin.Context) {
		pokemons := GetAllPokemons()

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Pokemons List",
			"data":    pokemons,
		})
	})

	router.GET("/pokemon/:name", func(c *gin.Context) {
		name := c.Param("name")

		reqStatus := http.StatusOK
		data := getPokemonByName(name)
		var message string
		if _, err := strconv.Atoi(name); err == nil {
			// check if name is number
			reqStatus = http.StatusBadRequest
			message = "Name must be a string"
			data = nil
		} else if data == nil && reqStatus == http.StatusOK {
			// check if pokemon exists
			reqStatus = http.StatusNotFound
			message = fmt.Sprintf("Pokemon '%s' not found", name)
		} else {
			// pokemon exists
			message = fmt.Sprintf("Pokemon %s", name)
		}

		c.JSON(reqStatus, gin.H{
			"status":  reqStatus,
			"message": message,
			"data":    data,
		})
	})

	if LOCAL_DEBUG {
		router.Run()
	} else {
		ginLambda = ginadapter.New(router)
		lambda.Start(handler)
	}
}

var ginLambda *ginadapter.GinLambda

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
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
