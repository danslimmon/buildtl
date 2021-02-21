package main

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// fsChangeSource returns a channel that will receive an empty struct on filesystem changes.
//
// dir is the directory in which to look (recursively, ignoring symlinks). done is a channel that
// should be closed when the change source is no longer being read from.
//
// This function returns a channel on which empty structs will be sent whenever there's a filesystem
// change in dir. An error is returned if there's a problem initializing the filesystem watcher. If
// the filesystem watcher closes its channel(s), the channel returned by fsChangeSource will be
// closed.
func fsChangeSource(dir string, done <-chan struct{}) (chan struct{}, error) {
	return nil, nil
}

// singleDirChangeSource returns a channel that will receive an empty struct on filesystem changes
// within a given directory.
func singleDirChangeSource(dir string, done <-chan struct{}) (chan struct{}, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Close watcher on done
	go func() {
		<-done
		watcher.Close()
	}()

	out := make(chan struct{})
	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				log.Info("received event")
				if !ok {
					close(out)
					return
				}
				out <- struct{}{}
			case err, ok := <-watcher.Errors:
				// Errors on this channel can be file removals, so treat them like any other event.
				if !ok {
					close(out)
					return
				}
				log.Error(err)
			case <-done:
				return
			}
		}
	}()

	if err := watcher.Add(dir); err != nil {
		return nil, err
	}

	return out, nil
}

func main() {
	return
}
