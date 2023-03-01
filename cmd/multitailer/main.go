package main

import (
	"context"
	"fmt"

	"github.com/nullren/multitailer/pkg/tailer"
)

func main() {
	tailer.Watch(context.Background(), ".", func(file, content string) error {
		fmt.Printf("%s: %s\n", file, content)
		return nil
	})
}
