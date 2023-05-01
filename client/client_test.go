package client

import (
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"github.com/pact-foundation/pact-go/dsl"
	"net/http"
	"os"
	"strings"
	"testing"
)

var pact dsl.Pact
var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/pacts", dir)
var logDir = fmt.Sprintf("%s/log", dir)

func createPact() dsl.Pact {
	return dsl.Pact{
		Consumer: "Sample Consumer",
		Provider: "Sample Provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	}
}

func TestMain(m *testing.M) {
	// Setup Pact and related test stuff
	pact = createPact()

	// Proactively start service to get access to the port
	pact.Setup(true)

	// Run all the tests
	code := m.Run()

	// Shutdown the Mock Service and write pact files to disk
	err := pact.WritePact()
	if err != nil {
		log.Infof("Failed to write your contract")
		return
	}

	pact.Teardown()
	os.Exit(code)
}

func getCar() (err error) {
	url := fmt.Sprintf("http://localhost:%d/cars/1", pact.Server.Port)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(req)

	return
}

type Car struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

var validCar = Car{ID: "30", Title: "Toyota", Color: "Yellow"}
var invalidCar = Car{Title: "Kia"}

func postCar() (err error) {
	url := fmt.Sprintf("http://localhost:%d/cars", pact.Server.Port)
	jsonPayload, _ := json.Marshal(validCar)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(req)

	return
}

func postInvalidCar() (err error) {
	url := fmt.Sprintf("http://localhost:%d/cars", pact.Server.Port)
	jsonPayload, _ := json.Marshal(invalidCar)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(req)

	return
}

func TestTheWholeBody_POST(t *testing.T) {
	pact.
		AddInteraction().
		Given("Validate the whole response body").
		UponReceiving("A POST request").
		WithRequest(dsl.Request{
			Method: "POST",
			Path:   dsl.Term("/cars", "/cars"),
			Body:   validCar,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 201,
			Body:   validCar,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(postCar)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func TestTheWholeBody_GET(t *testing.T) {
	pact.AddInteraction().
		Given("Validate the whole response body").
		UponReceiving("A GET request").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars/1", "/cars/[0-9]+"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   Car{ID: "1", Title: "BMW", Color: "Black"},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(getCar)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func TestSomeKeys_GET(t *testing.T) {
	pact.
		AddInteraction().
		Given("Validate title and color").
		UponReceiving("A GET request").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars/1", "/cars/[0-9]+"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body: map[string]interface{}{

				// Key + value validation
				"title": "BMW",

				// Validate exact key + the format of value
				"color": dsl.Term("Yellow", `\w+`),
			},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(getCar)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func TestSomeKeys_POST(t *testing.T) {
	pact.
		AddInteraction().
		Given("Validate title and color").
		UponReceiving("A POST request").
		WithRequest(dsl.Request{
			Method: "POST",
			Path:   dsl.Term("/cars", "/cars"),
			Body:   validCar,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 201,
			Body: map[string]interface{}{

				// Key + value validation
				"title": "Toyota",

				// Validate exact key + the format of value
				"color": dsl.Term("Yellow", `\w+`),
			},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(postCar)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func TestInvalidRequest_GET(t *testing.T) {
	pact.
		AddInteraction().
		Given("Validate error message").
		UponReceiving("A GET request with invalid ID").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars/9999", "/cars/[0-9]+"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 404,
			Body: map[string]string{
				"message": "Requested car is not found",
			},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(getCar)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func TestInvalidRequest_POST(t *testing.T) {
	pact.
		AddInteraction().
		Given("Validate error message").
		UponReceiving("A POST request with no ID provided").
		WithRequest(dsl.Request{
			Method: "POST",
			Path:   dsl.Term("/cars", "/cars"),
			Body:   invalidCar,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 400,
			Body: map[string]string{
				"message": "ID must not be empty",
			},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(postInvalidCar)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}
