package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

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
	case "limits":
		return runMessagesLimits(client, flag.Args()[1:])
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
		message    = fs.String("m", "", "Message to send")
		html       = fs.Bool("html", false, "Message is formatted as HTML")
		mono       = fs.Bool("mono", false, "Use monospace font for message")
		title      = fs.String("t", "", "Title (optional)")
		device     = fs.String("d", "", "Device (optional)")
		url        = fs.String("url", "", "URL (optional)")
		urlTitle   = fs.String("url-title", "", "URL title (optional)")
		priority   = fs.String("p", "", "Priority (lowest,low,normal,high, or emergency)")
		sound      = fs.String("sound", "", "Sound to play (optional)")
		attachment = fs.String("a", "", "File to attach (optional)")
		retry      = fs.Duration("retry", 30*time.Second, "Retry duration (if priority is emergency)")
		expire     = fs.Duration("expire", 5*time.Minute, "Expire duration (if priority is emergency)")
		tags       = fs.String("tags", "", "Tags (optional)")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	prio := pushover.Normal
	switch *priority {
	default:
		prio = pushover.Normal
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
		Message:    *message,
		HTML:       *html,
		Monospace:  *mono,
		Title:      *title,
		URL:        *url,
		URLTitle:   *urlTitle,
		Priority:   prio,
		Retry:      *retry,
		Expire:     *expire,
		Sound:      *sound,
		Attachment: *attachment,
	}
	if v := *device; v != "" {
		msg.Devices = []string{v}
	}
	if v := *tags; v != "" {
		msg.Tags = []string{v}
	}
	resp, err := client.Messages.Send(context.Background(), msg)
	if err != nil {
		return err
	}
	fmt.Println(resp.Receipt)

	return nil
}

func runMessagesLimits(client *pushover.Client, args []string) error {
	fs := flag.NewFlagSet("limits", flag.ExitOnError)
	fs.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s %s:\n", os.Args[0], fs.Name())
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	resp, err := client.Messages.Limits(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Limit=%d, remaining=%d, reset=%d (%s)\n",
		resp.Limit, resp.Remaining, resp.Reset,
		resp.ResetTime.Format(time.UnixDate))
	return nil
}
