package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type Bookmark struct {
	Title    string `yaml:"title"`
	Category string `yaml:"category"`
	Desc     string `yaml:"description"`
	Link     string `yaml:"link"`
}

//nolint: funlen // don't care about function length.
func main() {
	links := map[string]map[string]Bookmark{}

	if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yml") {
			// log.Printf("Process File: %s", path)
			yd := Bookmark{}

			yb, err := ioutil.ReadFile(path)
			if err != nil {
				log.Print(err)

				return fmt.Errorf("unable to read file(%s): %w", path, err)
			}
			if err := yaml.Unmarshal(yb, &yd); err != nil {
				log.Fatalf("error: %v", err)
			}
			if _, ok := links[yd.Category]; !ok {
				links[yd.Category] = make(map[string]Bookmark)
			}
			links[yd.Category][yd.Title] = yd
		}

		return nil
	}); err != nil {
		log.Print(err)

		return
	}

	f, err := os.Create("README.md")
	if err != nil {
		log.Print(err)

		return
	}
	defer f.Close()

	if _, err := fmt.Fprint(f, "# Bookmarks\n\n"); err != nil {
		log.Print(err)

		return
	}

	keys := []string{}
	for category := range links {
		keys = append(keys, category)
	}

	sort.Strings(keys)

	for _, category := range keys {
		if _, err := fmt.Fprintf(f, "## %s\n", category); err != nil {
			log.Print(err)

			return
		}

		ck := []string{}
		for title := range links[category] {
			ck = append(ck, title)
		}

		sort.Strings(ck)

		for _, title := range ck {
			if _, err := fmt.Fprintf(
				f,
				" - [%s](%s): %s\n",
				links[category][title].Title,
				links[category][title].Link,
				links[category][title].Desc,
			); err != nil {
				log.Print(err)

				return
			}
		}
	}
}
