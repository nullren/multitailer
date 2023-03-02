package main

import (
	"context"

	"github.com/nullren/multitailer/pkg/tailer"
)

func main() {
	if err := tailer.Watch(context.Background(), "/var/log", func(file, content string) error {
		// fmt.Printf("%s: %s\n", file, content)
		return nil
	}); err != nil {
		panic(err)
	}
}
