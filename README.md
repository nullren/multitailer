# multitailer

Multitailer is a Go library to tail and follow multiple files. A new tailer is
created by passing `NewMultitailer` a `MultitailerConfig`. The primary means of
using this libary is using the `Follow` function which an example of is in
`cmd/multitailer/main.go`.

```go
type MultitailerConfig struct {
	// CheckpointsSaveFile is the file to save file read checkpoints to.
	CheckpointsSaveFile string

	// CheckpointsSaveInterval is the interval to save checkpoints to disk.
	CheckpointsSaveInterval time.Duration

	// FileSearchGlob is the glob to search for files to watch.
	FileSearchGlob string

	// FileUpdateInterval is the interval to update the list of files to watch.
	FileUpdateInterval time.Duration

	// FileMaxReadBytes is the maximum number of bytes to read from a file
	FileMaxReadBytes int64
}

// FollowFunc is a function that is called when a file has new content. It is
// called for each new line data. Parameters are (file, line). Returning an
// error will end the outer Follow function calling this.
type FollowFunc = func(file, line string) error
```

New files are found periodically and any changes to existing files (moved or
truncated) will trigger re-reading the original file name from the start.

The last-read position is also kept and periodically written to disk.

## Notes

Inspiration came from originally reading [Vector.dev's kubernetes_logs
source](https://github.com/vectordotdev/vector/blob/ab459399a7ca58c088dfbd30dd6c08f5799c929e/src/sources/kubernetes_logs/mod.rs#L825-L838)
as it had a similar idea for checkpointing reads and periodically persisting
them to disk.
