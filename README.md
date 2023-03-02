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
