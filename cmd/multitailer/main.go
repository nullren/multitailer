package main

import (
	"context"
	"fmt"

	"github.com/nullren/multitailer/pkg/tailer"
)

func main() {
	if err := tailer.Watch(context.Background(), "/tmp/ren", func(file, line string) error {
		fmt.Printf("%s: %s\n", file, line)
		return nil
	}); err != nil {
		panic(err)
	}
}
