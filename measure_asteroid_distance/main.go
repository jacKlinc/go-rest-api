package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"net/http"
	// "strconv"
)

type Response events.APIGatewayProxyResponse

type Res struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	API       string `json:"api"`
}

// Todo struct
type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type NearEarthObjects struct {
	Links            string `json:"links"`
	Id               int    `json:"id"`
	NeoRefId         int    `json:"neo_reference_id"`
	Name             string `json:"name"`
	EstimateDiameter string `json:"estimated_diameter"`
}

type NasaResponse struct {
	Links           string `json:"links"`
	ElementCount    int    `json:"element_count"`
	NearEarthObject string `json:"near_earth_objects"`
}

var nasaRes NasaResponse

func Handler(ctx context.Context) (Response, error) {

	// 1. Call NASA REST API w/ params

	// 2. Parse reponse and put in Lambda format

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func nasaGet(start_date string, end_date string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	api := os.Getenv("NASA_API")

	NASA_URL := "https://api.nasa.gov/neo/rest/v1/feed"
	params := "?start_date=" + start_date + "&end_date=" + end_date + "&api_key=" + api

	fmt.Println("1. Performing Http Get...")
	resp, err := http.Get(NASA_URL + params)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println(resp.Body)

	// Convert response body to string
	bodyString := string(bodyBytes)
	// fmt.Println("API Response as String:\n" + bodyString)

	json.Unmarshal([]byte(bodyString), &nasaRes)
	fmt.Println("ElementCount: ", nasaRes.ElementCount)
	fmt.Println("NearEarthObjects: ", nasaRes.NearEarthObject)

	// out, _ := json.Marshal(bodyString)
	// fmt.Println("API Response as String:\n" + bodyString)

}

func main() {
	// lambda.Start(Handler)

	// res = append(res,
	// 	Res{
	// 		StartDate: "2015-09-07",
	// 		EndDate:   "2015-09-08",
	// 		API:       val,
	// 	},
	// )
	// fmt.Println(res)

	nasaGet("2015-09-07", "2015-09-08")

	// router := mux.NewRouter()

	// router.HandleFunc("/getAsteroids", getAsteroids).Methods("GET")

	// log.Fatal((http.ListenAndServe(":5000", router)))

}
