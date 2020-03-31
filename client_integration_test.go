// +build integration

package pushover_test

import (
	"context"
	"testing"

	"github.com/olivere/pushover-api-go"
)

func TestClientLifecycle(t *testing.T) {
	client, err := pushover.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	m := pushover.Message{
		Title:   "Pushover API for Go",
		Message: "This message is sent from a unit test of the Pushover API for Go at https://github.com/olivere/pushover-api-go",
	}
	if err := client.Messages.Send(context.Background(), m); err != nil {
		t.Fatalf("expected no error on sending message, got: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("expected client to close without error, got: %v", err)
	}
}
