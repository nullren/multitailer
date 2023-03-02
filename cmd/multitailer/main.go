package main

import (
	"context"
	"fmt"

	"github.com/nullren/multitailer"
)

func main() {
	if err := multitailer.Watch(context.Background(), "/tmp/ren", func(file, line string) error {
		fmt.Printf("%s: %s\n", file, line)
		return nil
	}); err != nil {
		panic(err)
	}
}
