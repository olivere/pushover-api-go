package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/olivere/pushover-api-go"
)

func runEnv(_ *pushover.Client, args []string) error {
	fs := flag.NewFlagSet("env", flag.ExitOnError)
	fs.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s %s:\n", os.Args[0], fs.Name())
		fmt.Fprintln(w)
		fmt.Fprint(w, "Defaults:\n")
		flag.PrintDefaults()
		fmt.Fprintln(w)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	fmt.Printf("%-30s: %s (required)\n", "APP_TOKEN", envString("", "APP_TOKEN"))
	fmt.Printf("%-30s: %s (required)\n", "USER_KEY", envString("", "USER_KEY"))
	fmt.Printf("%-30s: %s (required)\n", "PUSHOVER_URL", envString("", "PUSHOVER_URL"))

	return nil
}
