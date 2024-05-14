package onpremise

import (
	"fmt"
	"github.com/rjeczalik/notify"
	"path/filepath"
	"strings"
)

type fileWatcher interface {
	watch(path string, onChange func()) error
	unwatch(path string) error
	close() error
	run() error
}

type FileWatcher struct {
	logger         logWrapper
	callbacks      map[string]func()
	notifyChannels map[string]chan notify.EventInfo
}

func (f *FileWatcher) unwatch(path string) error {
	delete(f.callbacks, path)
	notify.Stop(f.notifyChannels[path])
	return nil
}

func (f *FileWatcher) run() error {
	for fileName, f2 := range f.callbacks {
		fileDirectory := filepath.Dir(fileName)
		eventChannel := make(chan notify.EventInfo, 10)
		err := notify.Watch(fileDirectory, eventChannel,
			notify.InMovedTo,
			notify.InCreate,
			notify.InModify,
		)
		if err != nil {
			return err
		}
		f2 := f2
		fileName := fileName
		go func() {
			for {
				select {
				case ei := <-eventChannel:
					fmt.Println(ei.Event().String(), ei.Path(), fileName)
					if strings.HasSuffix(ei.Path(), fileName) {
						f.logger.Printf("File %s has been modified", ei.Path())
						f2()
					}
				}
			}
		}()
		f.notifyChannels[fileName] = eventChannel
	}

	return nil
}

func (f *FileWatcher) watch(path string, onChange func()) error {
	//err := f.watcher.Add(path)
	//if err != nil {
	//	return err
	//}

	f.callbacks[path] = onChange

	return nil
}

func (f *FileWatcher) close() error {
	for _, c := range f.notifyChannels {
		notify.Stop(c)
	}

	return nil
}

func newFileWatcher(logger logWrapper) (*FileWatcher, error) {
	return &FileWatcher{
		logger:         logger,
		callbacks:      make(map[string]func()),
		notifyChannels: make(map[string]chan notify.EventInfo),
	}, nil
}
