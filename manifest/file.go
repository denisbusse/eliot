package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/layeryd/model"
)

// FileManifestSource is source what reads manifest from file
type FileManifestSource struct {
	filePath         string
	interval         time.Duration
	previousManifest []model.Pod
}

// NewFileManifestSource creates new file source what updates intervally
func NewFileManifestSource(filePath string, interval time.Duration) *FileManifestSource {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Panicf("Unable to open state, file [%s] does not exist!", filePath)
	}
	return &FileManifestSource{
		filePath: filePath,
		interval: interval,
	}
}

// GetUpdates return channel for manifest changes
func (s *FileManifestSource) GetUpdates() chan []model.Pod {
	updates := make(chan []model.Pod)
	go func() {
		for {
			time.Sleep(s.interval)

			pods, err := s.getPods()
			if err != nil {
				log.Printf("Error reading state: %s", err)
				continue
			}

			if reflect.DeepEqual(s.previousManifest, pods) {
				continue
			}

			updates <- pods
		}
	}()
	return updates
}

func (s *FileManifestSource) getPods() (pods []model.Pod, err error) {
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return pods, fmt.Errorf("Cannot update state, file [%s] does not exist", s.filePath)
		}
		return pods, err
	}

	switch extension := filepath.Ext(s.filePath); extension {
	case ".yaml", ".yml":
		return unmarshalYaml(data)
	default:
		return pods, fmt.Errorf("Invalid source file format: %s", extension)
	}
}