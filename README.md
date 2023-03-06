# multitailer

Multitailer is a Go library to tail and follow multiple files. A new tailer is
created by passing `NewMultitailer` a `MultitailerConfig`. The primary means of
using this libary is using the `Follow` function which an example of is in
`cmd/multitailer/main.go`.

New files are found periodically and any changes to existing files (moved or
truncated) will trigger re-reading the original file name from the start.

The last-read position is also kept and periodically written to disk.

## Notes

Inspiration came from originally reading [Vector.dev's kubernetes_logs
source](https://github.com/vectordotdev/vector/blob/ab459399a7ca58c088dfbd30dd6c08f5799c929e/src/sources/kubernetes_logs/mod.rs#L825-L838)
as it had a similar idea for checkpointing reads and periodically persisting
them to disk.
