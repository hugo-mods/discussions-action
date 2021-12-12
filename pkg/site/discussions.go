package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hugo-mods/discussions/bridge/pkg/model"
)

// Discussions maps site's blog post URL to Discussion.
type Discussions map[string]model.Discussion

func (d Discussions) Save(path string) error {
	data, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return fmt.Errorf("could not create directories to write discussions to JSON file: %v", err)
	}
	if err := os.WriteFile(path, data, 0666); err != nil {
		return fmt.Errorf("could not write discussions to JSON file: %v", err)
	}
	return nil
}

func (d Discussions) HasPage(url string) bool {
	_, ok := d[url]
	return ok
}

func (d Discussions) URLs() []string {
	urls := make([]string, len(d))
	for _, u := range urls {
		urls = append(urls, u)
	}
	return urls
}
