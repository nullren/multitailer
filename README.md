# multitailer

Multitailer is a Go library to tail and follow multiple files. The primary means of using
this libary is using the `Watch` function which an example of is in `cmd/multitailer/main.go`.

```go
// WatchFunc is a function that is called when a file has
// new content. It is called for each new line data.
// Parameters are (file, line). Returning an error will
// end the watch.
type WatchFunc = func(file, line string) error

// Watch watches a directory for files and new content and calls
// watchFunc for each new line of content.
func Watch(ctx context.Context, dir string, watchFunc WatchFunc) error
```

New files are found periodically and any changes to existing files (moved or truncated) will trigger re-reading the original file name from the start.

The last-read position is also kept and periodically written to disk.

## Notes

Inspiration came from originally reading [Vector.dev's kubernetes_logs source](https://github.com/vectordotdev/vector/blob/ab459399a7ca58c088dfbd30dd6c08f5799c929e/src/sources/kubernetes_logs/mod.rs#L825-L838) as it had a similar idea for checkpointing reads and periodically persisting them to disk.
