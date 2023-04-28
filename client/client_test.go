package client

import (
	"fmt"
	"github.com/google/martian/log"
	"github.com/pact-foundation/pact-go/dsl"
	"net/http"
	"os"
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

	// Shutdown the Mock Service and Write pact files to disk
	err := pact.WritePact()
	if err != nil {
		log.Infof("Failed to start you service")
		return
	}
	pact.Teardown()

	os.Exit(code)
}

func validateGet() (err error) {
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

func TestSomeValuesAndKeys_GET(t *testing.T) {
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
			Body: map[string]string{
				"title": "BMW",
			},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(validateGet)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func TestInvalidGetRequest_GET(t *testing.T) {
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

	err := pact.Verify(validateGet)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func TestAllKeys_GET(t *testing.T) {
	pact.
		AddInteraction().
		Given("Validate all keys are present").
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
			Body: map[string]dsl.Matcher{
				"title": dsl.Term("Toyota", `\w+`),
				"id":    dsl.Term("300", `\w+`),
				"color": dsl.Term("Yellow", `\w+`),
			},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(validateGet)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}

	// specify PACT publisher
	/*publisher := dsl.Publisher{}
	err = publisher.Publish(types.PublishRequest{
		PactURLs: []string{"../client/pacts/consumer_name-provider_name.json"},
		PactBroker:      "https://pen.pactflow.io/", //link to your remote Contract broker
		BrokerToken:     "jEQnxw7xWgYRv-3-G7Cx-g",   //your PactFlow token
		ConsumerVersion: "2.0.2",
		Tags:            []string{"2.0.2", "latest"},
	})*/
}

func TestTheWholeBody_GET(t *testing.T) {
	type Car struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Color string `json:"color"`
	}

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

	err := pact.Verify(validateGet)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}
