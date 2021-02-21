package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Channel returned by fsChangeSource should get message when a new file is created
func TestFSChangeSource_FileCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	done := make(chan struct{})
	defer close(done)

	// Make temp directory to work in
	tmpDir, err := ioutil.TempDir("", "timeline_test_*")
	defer os.RemoveAll(tmpDir)
	assert.Nil(err)

	// Start the goroutine that will listen for change events.
	src, err := fsChangeSource(tmpDir, done)
	assert.Nil(err)
	received := make(chan struct{})
	go func(received chan struct{}, done chan struct{}) {
		select {
		case <-src:
			close(received)
			return
		case <-done:
			return
		}
	}(received, done)

	// Create a file. This should trigger a signal to src.
	f, err := os.Create(filepath.Join(tmpDir, "foo"))
	defer f.Close()
	assert.Nil(err)

	// Wait for change event up to timeout.
	timeout := 100 * time.Millisecond
	select {
	case <-received:
		// Signal was received. Test passes.
		return
	case <-time.After(timeout):
		t.Logf("src did not receive message within timeout")
		t.FailNow()
	}
}

// Channel returned by fsChangeSource should get message when a file is modified
func TestFSChangeSource_FileModify(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	done := make(chan struct{})
	defer close(done)

	// Make temp directory to work in
	tmpDir, err := ioutil.TempDir("", "timeline_test_*")
	defer os.RemoveAll(tmpDir)
	assert.Nil(err)

	// Create a file, which we'll modify later to trigger a change signal.
	f, err := os.Create(filepath.Join(tmpDir, "foo"))
	defer f.Close()
	assert.Nil(err)

	// Start the goroutine that will listen for change events.
	src, err := fsChangeSource(tmpDir, done)
	assert.Nil(err)
	received := make(chan struct{})
	go func(received chan struct{}, done chan struct{}) {
		select {
		case <-src:
			close(received)
			return
		case <-done:
			return
		}
	}(received, done)

	// Modify foo. This should trigger a signal to src.
	_, err = f.WriteString("hello")
	assert.Nil(err)

	// Wait for change event up to timeout.
	timeout := 100 * time.Millisecond
	select {
	case <-received:
		// Signal was received. Test passes.
		return
	case <-time.After(timeout):
		t.Logf("src did not receive message within timeout")
		t.FailNow()
	}
}

// Channel returned by fsChangeSource should get message when a file is removed
func TestFSChangeSource_FileRemove(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	done := make(chan struct{})
	defer close(done)

	// Make temp directory to work in
	tmpDir, err := ioutil.TempDir("", "timeline_test_*")
	defer os.RemoveAll(tmpDir)
	assert.Nil(err)

	// Create a file, which we'll modify later to trigger a change signal.
	path := filepath.Join(tmpDir, "foo")
	f, err := os.Create(path)
	defer f.Close()
	assert.Nil(err)

	// Start the goroutine that will listen for change events.
	src, err := fsChangeSource(tmpDir, done)
	assert.Nil(err)
	received := make(chan struct{})
	go func(received chan struct{}, done chan struct{}) {
		select {
		case <-src:
			close(received)
			return
		case <-done:
			return
		}
	}(received, done)

	// Modify foo. This should trigger a signal to src.
	err = os.Remove(path)
	assert.Nil(err)

	// Wait for change event up to timeout.
	timeout := 100 * time.Millisecond
	select {
	case <-received:
		// Signal was received. Test passes.
		return
	case <-time.After(timeout):
		t.Logf("src did not receive message within timeout")
		t.FailNow()
	}
}
