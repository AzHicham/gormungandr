package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/CanalTP/gormungandr/internal/coverage"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func GetFileWithFS(uri url.URL) ([]*coverage.Coverage, error) {
	//Get all files in directory params
	log.Info("Mapping coverage-kraken, Read files from path :", uri.Path)
	fileInfo, err := ioutil.ReadDir(uri.Path)
	if err != nil {
		return nil, err
	}

	coverages := make([]*coverage.Coverage, 0)
	for _, file := range fileInfo {
		//filter to read only json files coverage
		if filepath.Ext(file.Name()) == ".json" {
			log.Info("Mapping coverage-kraken, Read file :", file.Name())
			f, err := os.Open(fmt.Sprintf("%s/%s", uri.Path, file.Name()))
			if err != nil {
				return nil, err
			}
			defer f.Close()
			logrus.Info("Read file: ", file.Name())
			var buffer bytes.Buffer
			if _, err = buffer.ReadFrom(f); err != nil {
				return nil, err
			}
			jsonData, err := ioutil.ReadAll(&buffer)
			if err != nil {
				return nil, err
			}
			coverage := &coverage.Coverage{}
			err = json.Unmarshal([]byte(jsonData), coverage)
			if err != nil {
				return nil, err
			}
			coverages = append(coverages, coverage)
		}
	}
	return coverages, nil
}
