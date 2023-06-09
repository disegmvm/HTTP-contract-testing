package main

import (
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"os"
	"testing"
)

func startProvider() {
	main()
}

func TestServerPact_Verification(t *testing.T) {

	go startProvider()

	var dir, _ = os.Getwd()
	var pactDir = fmt.Sprintf("%s/../client/pacts", dir)
	var logDir = fmt.Sprintf("%s/log", dir)

	pact := &dsl.Pact{
		Provider:                 "Sample Provider",
		LogDir:                   logDir,
		PactDir:                  pactDir,
		DisableToolValidityCheck: true,
	}

	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:            "http://127.0.0.1:8080",   //provider's URL
		BrokerURL:                  "https://pen.pactflow.io", //link to your remote Contract broker
		BrokerToken:                "jEQnxw7xWgYRv-3-G7Cx-g",  //your PactFlow token
		PublishVerificationResults: true,
		ProviderVersion:            "1.0.0",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Uncomment to verify contract locally
	/*log.Println("[debug] start verification")
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: "http://127.0.0.1:8080",
		PactURLs:        []string{"../client/pacts/sample_consumer-sample_provider.json"},
	})
	if err != nil {
		t.Fatal(err)
	}*/
}
