package filewatcher

import (
	"testing"
	"context"
	"time"
)


func TestNewFileWatcher(t *testing.T) {
	fw := NewFileWatcher()

	if fw == nil {
		t.Fatal("FW should not be nil")
	}

	if fw.executionInterval == 0 {
		t.Fatal("Should not be 0 as executionInterval")
	}
	t.Log(fw.executionInterval)
}

func TestWatchDirectory(t *testing.T) {
	 fw := NewFileWatcher()

	 ctx, cancel := context.WithCancel(context.TODO())
	fw.cancelBuffer = cancel
	 output := make(chan string)

	 fw.WatchDirectory(ctx, output, "testfiles/")

	 cancelTick := time.NewTicker(3 * time.Second)
	 for {
	 	select {
	 		case err := <- fw.ErrorChan:
	 			t.Log(err)
	 		case file := <- output:
	 			t.Log("Found file: ", file)
	 			case <- cancelTick.C:
	 				return
		}
	 }
}