package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/olivere/pushover-api-go"
)

var (
	insecure = flag.Bool("k", false, "Accept insecure connections")
	verbose  = flag.Bool("v", false, "Verbose output to stderr")
	raw      = flag.Bool("d", false, "Raw output of HTTP request/response to stderr")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
	log.SetOutput(os.Stdout)

	if err := runMain(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func runMain() error {
	rand.Seed(time.Now().UnixNano())
	flag.Usage = usage
	flag.Parse()

	// Set up client
	options := []pushover.ClientOption{}
	if *insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		options = append(options, pushover.WithTransport(tr))
	}
	if *verbose {
		options = append(options, pushover.WithLogger(pushover.NewJSONLogger(os.Stderr)))
	} else if *raw {
		options = append(options, pushover.WithLogger(pushover.NewRawLogger(os.Stderr)))
	}
	client, err := pushover.NewClient(options...)
	if err != nil {
		return err
	}
	defer client.Close()

	switch flag.Arg(0) {
	default:
		flag.Usage()
		os.Exit(2)
	case "env":
		return runEnv(client, flag.Args()[1:])
	case "messages":
		return runMessages(client, flag.Args()[1:])
	case "send":
		return runMessagesSend(client, flag.Args()[1:])
	}
	return nil
}

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
	fmt.Fprintln(w)
	fmt.Fprint(w, "Global defaults:\n")
	flag.PrintDefaults()
	fmt.Fprintln(w)
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "  env          Print environment\n")
	fmt.Fprint(w, "  send         Send a message\n")
	fmt.Fprintln(w)
}

func envString(defaultValue string, keys ...string) string {
	for _, key := range keys {
		if s := os.Getenv(key); s != "" {
			return s
		}
	}
	return defaultValue
}
