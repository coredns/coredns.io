package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	pluginDir = ".coredns/plugin/"
	pageDir   = "content/plugins/"
)

type page struct {
	header  meta
	content string
}

type meta struct {
	Title       string   `toml:"title"`
	Description string   `toml:"description"`
	Weight      int      `toml:"weight"`
	Tags        []string `toml:"tags"`
	Categories  []string `toml:"categories"`
	Date        string   `toml:"date"`
}

func main() {
	// Current pages
	pages, err := getPages(pageDir)
	if err != nil {
		log.Fatalf("Failed getting all page contents: %s", err)
	}

	// Incoming content
	plugins, err := getContent(pluginDir)
	if err != nil {
		log.Fatalf("Failed getting all plugin contents: %s", err)
	}

	neg, pos := compare(plugins, pages)

	if neg == nil && pos == nil {
		log.Printf("Nothing to update.")
	}

	// Remove any plugin files not present upstream
	for _, p := range neg {
		path := path.Join(pageDir, p)
		err := os.Remove(path + ".md")
		if err != nil {
			log.Fatalf("Error deleting file file %s: %s", path, err)
		}
		log.Printf("Removed content file: %s", path+".md")
	}

	// Overwrite/create necessary content files with toml header
	for _, p := range pos {
		h := plugins[p].header
		h.Date = time.Now().UTC().Format("2006-01-02T15:04:05.877581")
		h.Tags = []string{"plugin", h.Title}
		h.Categories = []string{"plugin"}

		var buf bytes.Buffer
		err := toml.NewEncoder(&buf).Encode(h)
		if err != nil {
			log.Fatalf("Error creating TOML header data: %s", err)
		}

		headerData := buf.String()
		contentData := plugins[p].content
		filePath := path.Join(pageDir, p)
		fileContent := strings.Join([]string{"", headerData, contentData}, "+++\n")
		err = ioutil.WriteFile(filePath+".md", []byte(fileContent), 0644)
		if err != nil {
			log.Fatalf("Error writing file: %s", err)
		}
		log.Printf("Wrote content changes %s to %s", path.Join(pluginDir, p, "README.md"), filePath)
	}
}

// Get current content pages from hugo content directory
// Returned map doesn't include header data as this isn't needed currently
func getPages(pageDir string) (pages map[string]page, err error) {
	pages = make(map[string]page)
	pd, err := ioutil.ReadDir(pageDir)
	if err != nil {
		return nil, err
	}

	for _, p := range pd {
		if p.IsDir() {
			continue
		}

		path := path.Join(pageDir, p.Name())
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var content string
		headerSeparator := 0

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if headerSeparator == 2 {
				content += scanner.Text() + "\n"
			}

			if scanner.Text() == "+++" {
				headerSeparator++
			}
		}
		pageData := page{content: content}
		pages[strings.Replace(p.Name(), ".md", "", -1)] = pageData
	}
	return pages, nil
}

// Get latest content from coreDNS code directory
func getContent(pluginDir string) (plugins map[string]page, err error) {
	plugins = make(map[string]page)
	pl, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		return nil, err
	}

	weight := 0

	for _, p := range pl {
		if !p.IsDir() {
			continue
		}
		if getBlacklist()[p.Name()] {
			continue
		}
		weight++

		path := path.Join(pluginDir, p.Name(), "README.md")
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var header meta
		var content string
		headerStarted := false
		contentStarted := false

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if contentStarted {
				content += scanner.Text() + "\n"
			}

			if headerStarted && scanner.Text() != "" {
				rawHeader := strings.Split(scanner.Text(), " - ")
				title := strings.Replace(rawHeader[0], "*", "", -1)
				header.Title = title
				header.Description = strings.Join([]string{"*", title, "* ", rawHeader[1]}, "")
				header.Weight = weight
				headerStarted = false
				contentStarted = true
			}

			if scanner.Text() == "## Name" {
				headerStarted = true
			}
		}
		pluginData := page{header: header, content: content}
		plugins[p.Name()] = pluginData
	}
	return plugins, nil
}

// Compare current and incoming state
// Return a negative list to delete
// Return a positive list to create/overwrite
func compare(plugins map[string]page, pages map[string]page) (neg []string, pos []string) {
	// Prepare list of pages to delete
	for k := range pages {
		if _, ok := plugins[k]; !ok {
			neg = append(neg, k)
		}
	}

	// Prepare list of pages to create
	for k := range plugins {
		// New pages to create
		if _, ok := pages[k]; !ok {
			pos = append(pos, k)
		}

		// Check necessity for an update
		if plugins[k].content != pages[k].content {
			pos = append(pos, k)
		}
	}
	return neg, pos
}

// Return directory blacklist
func getBlacklist() map[string]bool {
	var bl = map[string]bool{
		"test":       true,
		"pkg":        true,
		"deprecated": true,
	}
	return bl
}
