package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/olivere/pushover-api-go"
)

func runMessages(client *pushover.Client, args []string) error {
	fs := flag.NewFlagSet("messages", flag.ExitOnError)
	fs.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s %s:\n", os.Args[0], fs.Name())
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	switch flag.Arg(1) {
	default:
		return fmt.Errorf("unsupported call: %s", flag.Arg(1))
	case "send":
		return runMessagesSend(client, flag.Args()[1:])
	}
}

func runMessagesSend(client *pushover.Client, args []string) error {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	fs.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s %s:\n", os.Args[0], fs.Name())
		fs.PrintDefaults()
	}
	var (
		message  = fs.String("m", "", "Message to send")
		html     = fs.Bool("html", false, "Message is formatted as HTML")
		mono     = fs.Bool("mono", false, "Use monospace font for message")
		title    = fs.String("t", "", "Title (optional)")
		device   = fs.String("d", "", "Device (optional)")
		url      = fs.String("url", "", "URL (optional)")
		urlTitle = fs.String("url-title", "", "URL title (optional)")
		priority = fs.String("p", "", "Priority (lowest,low,normal,high, or emergency)")
		sound    = fs.String("sound", "", "Sound to play (optional)")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	prio := pushover.Normal
	switch *priority {
	default:
		return fmt.Errorf("unknown priority: %s", *priority)
	case "lowest":
		prio = pushover.Lowest
	case "low":
		prio = pushover.Low
	case "normal":
		prio = pushover.Normal
	case "high":
		prio = pushover.High
	case "emergency":
		prio = pushover.Emergency
	}

	msg := pushover.Message{
		Message:   *message,
		HTML:      *html,
		Monospace: *mono,
		Title:     *title,
		Devices:   []string{*device},
		URL:       *url,
		URLTitle:  *urlTitle,
		Priority:  prio,
		Sound:     *sound,
	}
	resp, err := client.Messages.Send(context.Background(), msg)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)

	return nil
}
