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

type Car struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

var createCar = Car{ID: "30", Title: "Toyota", Color: "Yellow"}

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

	// Shutdown the Mock Service and Write pact files to disk
	err := pact.WritePact()
	if err != nil {
		log.Infof("Failed to start you service")
		return
	}
	pact.Teardown()

	os.Exit(code)
}

func getCar() (err error) {
	url := fmt.Sprintf("http://localhost:%d/cars/1", pact.Server.Port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	return
}

func postCar() (err error) {
	url := fmt.Sprintf("http://localhost:%d/cars", pact.Server.Port)
	jsonPayload, _ := json.Marshal(createCar)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	return
}

func postInvalidCar() (err error) {
	url := fmt.Sprintf("http://localhost:%d/cars", pact.Server.Port)
	jsonPayload, _ := json.Marshal(createCar.ID)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	return
}

func TestTheWholeBody_POST(t *testing.T) {
	pact.
		AddInteraction().
		Given("Validate the whole body").
		UponReceiving("A POST request").
		WithRequest(dsl.Request{
			Method: "POST",
			Path:   dsl.Term("/cars", "/cars"),
			Body:   createCar,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 201,
			Body:   createCar,
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
		Given("Match the whole response body").
		UponReceiving("A a GET request").
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
		Given("Validate title only").
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
				"title": "BMW",
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
		Given("Validate title only").
		UponReceiving("A POST request").
		WithRequest(dsl.Request{
			Method: "POST",
			Path:   dsl.Term("/cars", "/cars"),
			Body:   createCar,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 201,
			Body: map[string]interface{}{
				"title": "Toyota",
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
		Given("Error message").
		UponReceiving("A GET invalid request").
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
		Given("Error message").
		UponReceiving("A POST request").
		WithRequest(dsl.Request{
			Method: "POST",
			Path:   dsl.Term("/cars", "/cars"),
			Body:   createCar.ID,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 400,
			//TODO response is in double quotes
			//Body:   "Provided data is invalid",
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("text/plain; charset=utf-8", `text\/plain`),
			},
		})

	err := pact.Verify(postInvalidCar)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}
