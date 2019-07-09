fileWatcher package is a utility to monitor new files in a directory.
It's mainly used as an ingress monitor for file based transfers.


FileWatcher has a TTL which is used by a buffer to avoid reporting new files.
Those files are stored in the buffered in TTL (seconds)
When time has ran out, they will be reported again, it is up to the using system to remove / move files.

ExecutionTime is how often the directory should be monitored.

Recommended usage:

	// Create a channel with string
	filechannel := make(chan string)
	// Create a file watcher
	watcher := filewatcher.NewFileWatcher()
	watcher.ChangeExecutionTime(1)
	watcher.ChangeTTL(5)
	go watcher.WatchDirectory(filechannel, "./fileWatcher/")

	for {
		select {
		case newFile := <-filechannel:
			fmt.Println(newFile)
		}
	}
