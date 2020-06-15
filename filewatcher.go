// Package filewatcher is used to monitor directories and their activity for new files
// @version 1.0
// @author Percy Bolmer
package filewatcher

import (
	"io/ioutil"
	"log"
	"sync"
	"context"
	"time"
)

// executionInterval is how often we should monitor that specific target
var executionIntervaldefault = time.Second * 10

// ttl is used to describe how long the files can live in memory
var ttldeafult int64 = 3600

// FileWatcher is used to monitor a directory or files for events
type FileWatcher struct {
	sync.Mutex
	found map[string]int64
	ttl int64
	executionInterval time.Duration

	cancelBuffer context.CancelFunc
}

// ChangeTTL is used to change how long the files should live in the buffer memory
// The files wont be viewed as new files until the TTL has ran out, this is to avoid memory leaks
// ttl should be given in seconds, and default is 3600
func (watcher *FileWatcher) ChangeTTL(TTL int64) {
	watcher.ttl = TTL
	return
}

// ChangeExecutionTime is to set how often the directory should be monitored
// interval should be given in seconds between each execution
func (watcher *FileWatcher) ChangeExecutionTime(interval int) {
	watcher.executionInterval = time.Second * time.Duration(interval)
	return
}

func (watcher *FileWatcher) Stop() {
	watcher.cancelBuffer()

}

// NewFileWatcher create an new words object
func NewFileWatcher() (watcher *FileWatcher) {
	watcher = &FileWatcher{found: map[string]int64{}}
	// Starts a gofunction that checks the timestamp on the items, remove them if neccessary
	ctxroot := context.Background()
	ctx, cancel := context.WithCancel(ctxroot)
	watcher.cancelBuffer = cancel
	go func() {
		interval := time.Tick(time.Second * 1)
		for {
			select {
				case now := <- interval:
					watcher.Lock()
					for k, v := range watcher.found {
						if now.Unix()-v > watcher.ttl {
							delete(watcher.found, k) // If the item is older than given time setting, delete it from buffer
						}
					}
					watcher.Unlock()
					case <- ctx.Done():
						return
			}
		}
	}()
	return
}

// WatchDirectory monitors a directory and returns files and new files to the given Channel
// It monitors for the set default time
func (watcher *FileWatcher) WatchDirectory(ctx context.Context, out chan<- string, directoryPath string) {
	go func(){
		ticker := time.NewTicker(watcher.executionInterval * time.Second)
		for {
			select {
				case <- ticker.C:
					files, err := ioutil.ReadDir(directoryPath)
					if err != nil {
						log.Fatal(err)
					}
					for _, f := range files {
						if f.IsDir() == false {
							// Add the File to the found map so that we dont send the same on back again
							_, ok := watcher.found[f.Name()]
							if !ok {
								watcher.Lock()
								watcher.found[f.Name()] = time.Now().Unix()
								watcher.Unlock()
								out <- f.Name()
							}

						}
					}
				case <- ctx.Done():
					return
			}
		}
	}()
}
