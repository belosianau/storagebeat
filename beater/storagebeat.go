package beater

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"storagebeat/config"
)

type Storagebeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

type StorageInfo struct {
	filename string `json:"storage.filename"`
	size     int64  `json:"storage.size"`
}

// New creates an instance of storagebeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Storagebeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts storagebeat.
func (bt *Storagebeat) Run(b *beat.Beat) error {
	bt.client, _ = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		storageInfos := bt.collectStorageStats(bt.config.Filesets)

		for _, storageInfo := range *storageInfos {
			bt.client.Publish(makeEvent(&storageInfo))
		}
	}
}

func (bt *Storagebeat) collectStorageStats(filesets *[]string) *[]StorageInfo {
	storageInfos := make([]StorageInfo, len(*filesets))
	for index, filename := range *filesets {
		storageInfos[index] = getFileSize(filename)
	}
	return &storageInfos
}

func getFileSize(filename string) StorageInfo {
	storageInfo := StorageInfo{}
	storageInfo.filename = filename

	fi, err := os.Stat(filename)
	if err != nil {
		return storageInfo
	}

	switch mode := fi.Mode(); {
	case mode.IsRegular():
		storageInfo.size = fi.Size()
	}

	return storageInfo
}

func makeEvent(storageInfo *StorageInfo) beat.Event {
	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":     "storagebeat",
			"filename": storageInfo.filename,
			"size":     storageInfo.size,
		},
	}
	return event
}

// Stop stops storagebeat.
func (bt *Storagebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
