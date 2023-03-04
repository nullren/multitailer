package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/nullren/multitailer"
)

func main() {
	searchGlob := flag.String("s", "/var/log/pods/*/*/*.log", "the search glob for files to tail")
	flag.Parse()

	watch, err := multitailer.NewWatch(multitailer.WatchConfig{
		CheckpointsSaveFile:     "/tmp/checkpoints.json",
		CheckpointsSaveInterval: 5 * time.Second,
		FileSearchGlob:          *searchGlob,
		FileUpdateInterval:      5 * time.Second,
		FileMaxReadBytes:        1024,
	})
	if err != nil {
		panic(err)
	}

	err = watch.Watch(context.Background(), func(file, line string) error {
		fmt.Printf("%s: %s\n", file, line)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
