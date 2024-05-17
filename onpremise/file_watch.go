package onpremise

import (
	"time"
)

type fileWatcher interface {
	watch(onChange func()) error
	stop() error
	run() error
}

type FileWatcher struct {
	watcher  *Watcher
	logger   logWrapper
	callback func()
	stopCh   chan struct{}
}

func (f *FileWatcher) run() error {
	for {
		select {
		case <-f.stopCh:
			f.stop()
			return nil
		case event := <-f.watcher.Changed():
			f.logger.Printf("File %s has been modified", event.Name())
			f.callback()
		case err, ok := <-f.watcher.Errors():
			if !ok {
				return err
			}
			if err != nil {
				f.logger.Printf("Error watching file: %v", err)
			}
		}
	}
}

func (f *FileWatcher) watch(onChange func()) error {
	f.callback = onChange
	return f.watcher.Start()
}

func (f *FileWatcher) stop() error {
	return f.watcher.Stop()
}

func newFileWatcher(logger logWrapper, path string, stopCh chan struct{}) (*FileWatcher, error) {
	watcher, err := newWatcher(path, time.Millisecond, time.Second)
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		watcher:  watcher,
		logger:   logger,
		callback: func() {},
		stopCh:   stopCh,
	}, nil
}
