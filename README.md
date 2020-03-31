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

## Command-line client

There is a simple command line client included to illustrate the usage
of the client. To build it, run:

```sh
go build ./cmd/pushover
```

To get a list of all supported commands and sub-commands, run:

```sh
$ ./pushover help
Usage of ./pushover:

Global defaults:
  -d	Raw output of HTTP request/response to stderr
  -k	Accept insecure connections
  -v	Verbose output to stderr

Commands:
  env          Print environment
  send         Send a message
```

Here's an example of how to send a message with an attachment:

```sh
export APP_TOKEN=...
export USER_KEY=...
./pushover send -t "Introduction" -m "Here's an avatar of mine." -a ~/Pictures/Avatar.png
```

To get a list of all options of a command, add `-h` to the command, e.g.:

```sh
$ ./pushover send -h
Usage of ./pushover send:
  -a string
    	File to attach (optional)
  -d string
    	Device (optional)
  -expire duration
    	Expire duration (if priority is emergency) (default 5m0s)
  -html
    	Message is formatted as HTML
  -m string
    	Message to send
  -mono
    	Use monospace font for message
  -p string
    	Priority (lowest,low,normal,high, or emergency)
  -retry duration
    	Retry duration (if priority is emergency) (default 30s)
  -sound string
    	Sound to play (optional)
  -t string
    	Title (optional)
  -tags string
    	Tags (optional)
  -url string
    	URL (optional)
  -url-title string
    	URL title (optional)
```

## License

MIT.
