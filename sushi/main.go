package main

// https://medium.com/codezillas/building-a-restful-api-with-go-part-1-9e234774b14d

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
)

type Roll struct {
	ID          string `json:"id"`
	ImageNumber string `json:"imageNumber"`
	Name        string `json:"name"`
	Ingredients string `json:"ingredients"`
}

var rolls []Roll

func getRolls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rolls)
}

func getRoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range rolls {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createRoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newRoll Roll
	json.NewDecoder(r.Body).Decode(&newRoll)
	newRoll.ID = strconv.Itoa(len(rolls) + 1)
	rolls = append(rolls, newRoll)
	json.NewEncoder(w).Encode(newRoll)
}

func updateRoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range rolls {
		if item.ID == params["id"] {
			rolls = append(rolls[:i], rolls[i+1:]...)
			var newRoll Roll
			json.NewDecoder(r.Body).Decode(&newRoll)
			newRoll.ID = params["id"]
			json.NewEncoder(w).Encode(newRoll)
			return
		}
	}
}

func deleteRoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range rolls {
		if item.ID == params["id"] {
			rolls = append(rolls[i:], rolls[i+1:]...)
			break
		}
	}
}

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {

	rolls = append(rolls,
		Roll{
			ID:          "1",
			ImageNumber: "8",
			Name:        "Spicy Tuna Roll",
			Ingredients: "Tuna, Chili sauce, Nori, Rice",
		}, Roll{
			ID:          "2",
			ImageNumber: "6",
			Name:        "California Roll",
			Ingredients: "Crab, Avocado, Cucumber, Nori, Rice",
		})

	// router := mux.NewRouter()

	// endpoints
	// router.HandleFunc("/sushi", getRolls).Methods("GET")
	// router.HandleFunc("/sushi", getRoll).Methods("GET")
	// router.HandleFunc("/sushi", createRoll).Methods("POST")
	// router.HandleFunc("/sushi", updateRoll).Methods("POST")
	// router.HandleFunc("/sushi", deleteRoll).Methods("DELETE")

	// log.Fatal((http.ListenAndServe(":5000", router)))

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

func main() {
	lambda.Start(Handler)
}
