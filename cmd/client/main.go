package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Meh-Mehul/sf/client"
)

func main() {
	hubFlag := flag.String("hub", "localhost:9999", "central hub address")
	flag.Parse()

	opts := &client.ClientOpts{
		Address: *hubFlag,
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Connected to %s\n", *hubFlag)
	fmt.Println("Type 'exit' to quit.")

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command := strings.TrimSpace(line)

		if command == "" {
			continue
		}

		if command == "exit" {
			break
		}

		out, code, err := client.PushJob(opts, command)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		fmt.Print(out)

		if code != 0 {
			fmt.Printf("[exit code: %d]\n", code)
		}
	}
}
