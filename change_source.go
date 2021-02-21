package main

import (
	"io/fs"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// singleDirChangeSource returns a channel that will receive an empty struct on filesystem changes
// within a given directory.
//
// dir is the directory in which to look (non-recursively). done is a channel that should be closed
// when the change source is no longer being read from.
//
// This function returns a channel on which empty structs will be sent whenever there's a filesystem
// change in dir. An error is returned if there's a problem initializing the filesystem watcher. If
// the filesystem watcher closes its channel(s), the channel returned by fsChangeSource will be
// closed.
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
	// The channel we'll return, which will contain the merged output from all the source channels
	// created inside filepath.Walk.
	out := make(chan struct{})
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		src, err := singleDirChangeSource(path, done)
		if err != nil {
			return err
		}
		go func() {
			for {
				select {
				case <-src:
					out <- struct{}{}
				case <-done:
					close(out)
					return
				}
			}
		}()
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func main() {
	return
}
