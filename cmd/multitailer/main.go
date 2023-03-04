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

	mt, err := multitailer.NewMultitailer(multitailer.MultitailerConfig{
		CheckpointsSaveFile:     "/tmp/checkpoints.json",
		CheckpointsSaveInterval: 5 * time.Second,
		FileSearchGlob:          *searchGlob,
		FileUpdateInterval:      5 * time.Second,
		FileMaxReadBytes:        10 * 1024 * 1024,
		FileLoopPauseTime:       1 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	err = mt.Follow(context.Background(), func(file, line string) error {
		fmt.Printf("%s: %s\n", file, line)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
