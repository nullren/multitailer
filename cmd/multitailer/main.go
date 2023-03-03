package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/nullren/multitailer"
)

func main() {
	searchGlob := flag.String("s", "/var/log/pods/*/*/*.log", "the search glob for files to tail")
	flag.Parse()

	if err := multitailer.Watch(context.Background(), *searchGlob, func(file, line string) error {
		fmt.Printf("%s: %s\n", file, line)
		return nil
	}); err != nil {
		panic(err)
	}
}
