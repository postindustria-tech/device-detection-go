package onpremise

import (
	"github.com/fsnotify/fsnotify"
)

type fileWatcher interface {
	watch(path string, onChange func()) error
	unwatch(path string) error
	close() error
	run() error
}

type FileWatcher struct {
	watcher   *fsnotify.Watcher
	logger    logWrapper
	callbacks map[string]func()
}

func (f *FileWatcher) unwatch(path string) error {
	return f.watcher.Remove(path)
}

func (f *FileWatcher) run() error {
	for {
		select {
		case event := <-f.watcher.Events:
			if event.Has(fsnotify.Write) {
				f.logger.Printf("File %s has been modified", event.Name)
				if callback, ok := f.callbacks[event.Name]; ok {
					callback()
				}
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return err
			}
			if err != nil {
				f.logger.Printf("Error watching file: %v", err)
			}
		}
	}
}

func (f *FileWatcher) watch(path string, onChange func()) error {
	err := f.watcher.Add(path)
	if err != nil {
		return err
	}

	f.callbacks[path] = onChange

	return nil
}

func (f *FileWatcher) close() error {
	return f.watcher.Close()
}

func newFileWatcher(logger logWrapper) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		watcher:   watcher,
		logger:    logger,
		callbacks: make(map[string]func()),
	}, nil
}
