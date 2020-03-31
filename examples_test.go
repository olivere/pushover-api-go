package pushover_test

import (
	"context"
	"os"

	"github.com/olivere/pushover-api-go"
)

func Example() {
	// Configure and create the client
	client, err := pushover.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Run commands on the client
	resp, err := client.Messages.Send(context.Background(), pushover.Message{
		Message: "Hello world!",
	})
	if err != nil {
		panic(err)
	}
	_ = resp
}

func ExampleNewClient_default() {
	// NewClient will use these environment variables to construct the configuration
	// - PUSHOVER_URL
	// - APP_TOKEN
	// - USER_KEY
	client, err := pushover.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
}

func ExampleNewClient_withTokenAndUser() {
	client, err := pushover.NewClient(
		pushover.WithAppToken("..."),
		pushover.WithUserKey("..."),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()
}

func ExampleNewClient_withLogger() {
	// Log HTTP request and response to stderr
	client, err := pushover.NewClient(
		pushover.WithLogger(pushover.NewRawLogger(os.Stderr)),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()
}
