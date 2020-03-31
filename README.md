# Pushover API for Go

[![Actions Status](https://github.com/olivere/pushover-api-go/workflows/Test/badge.svg)](https://github.com/olivere/pushover-api-go/actions)

This library allows programmatic access to the Pushover API with Go.

## In a nutshell

Here are a few code snippets of how to use the Pushover API with Go.

### Create a client

```go
import (
    "github.com/olivere/pushover-api-go"
)

client, err := pushover.NewClient(
    pushover.WithAppToken("..."),
    pushover.WithUserKey("..."),
)
if err != nil {
    return fmt.Errorf("unable to create client: %w", err)
}
defer client.Close()
```

### Send a message

```go
msg := pushover.Message{
    Title:   "Invitation",
    Message: "Please come over to ...",
}
err := client.Messages.Send(context.Background(), msg)
if err != nil {
    return err
}
```

## License

MIT.
