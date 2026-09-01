package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Meh-Mehul/sf/client"
)

func main() {
	hubFlag := flag.String("hub", "localhost:9999", "central hub address")
	cmdFlag := flag.String("cmd", "", "command to execute on remote server")
	flag.Parse()

	targetCmd := *cmdFlag
	if targetCmd == "" && len(flag.Args()) > 0 {
		targetCmd = strings.Join(flag.Args(), " ")
	}

	if targetCmd == "" {
		fmt.Fprintln(os.Stderr, "Usage: client [-hub <addr>] -cmd \"<command>\" OR client [-hub <addr>] <command...>")
		os.Exit(1)
	}

	opts := &client.ClientOpts{
		Address: *hubFlag,
	}

	out, code, err := client.PushJob(opts, targetCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(code)
	}

	fmt.Print(out)
	os.Exit(code)
}
