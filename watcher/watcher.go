package watcher

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type FileWatcher struct {
	v  []*viper.Viper
	mu sync.Mutex
}

func New() *FileWatcher {
	watcher := &FileWatcher{}
	return watcher
}

func (fw *FileWatcher) Add(filename string, fn func(in fsnotify.Event)) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	v := viper.New()
	v.SetConfigFile(filename)
	v.WatchConfig()
	v.OnConfigChange(fn)
	fw.v = append(fw.v, v)
}
