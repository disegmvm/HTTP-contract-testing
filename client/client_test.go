package client

import (
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"net/http"
	"os"
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
	pact.WritePact()
	pact.Teardown()

	os.Exit(code)
}

func Test_ContractTest(t *testing.T) {
	var test = func() (err error) {
		url := fmt.Sprintf("http://localhost:%d/cars", pact.Server.Port)
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

	type car struct {
		ID    dsl.Matcher `json:"id" pact:"example=1"`
		Title dsl.Matcher `json:"title" pact:"example=VAZ"`
		Color dsl.Matcher `json:"color" pact:"example=Nothing"`
	}

	type Carz struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Color string `json:"color"`
	}

	// uncomment to fail the validation
	//var zzz = Carz{
	//	ID: "1234e2", Title: "BMcvdfW", Color: "Bldcdcack",
	//}

	pact.
		AddInteraction().
		Given("Valid request").
		UponReceiving("A GET request to provider").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars", "/cars"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body: car{
				ID:    dsl.Like("1"),
				Title: dsl.Like("Vads"),
				Color: dsl.Like("Gjhtd"),
			},

			// uncomment to fail the validation
			//Body: zzz,

			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(test)
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
	})
	if err != nil {
		t.Fatal(err)
	}*/

}

func TestTheWholeResponseBody(t *testing.T) {
	var test = func() (err error) {
		url := fmt.Sprintf("http://localhost:%d/cars", pact.Server.Port)
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

	cars := []Car{
		{ID: "1", Title: "BMW", Color: "Black"},
		{ID: "2", Title: "Tesla", Color: "Red"},
	}

	pact.AddInteraction().
		Given("Validate the whole response body").
		UponReceiving("A a GET request").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars", "/cars"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   dsl.Match(cars),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(test)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}
