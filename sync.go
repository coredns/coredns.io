package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
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
	Title       string
	Description string
	Weight      int
	Tags        []string
	Categories  []string
	Date        time.Time
}

func main() {
	// Current pages
	pages := make(map[string]page)

	pd, err := ioutil.ReadDir(pageDir)
	if err != nil {
		log.Fatalf("Failed to get list of pages: %s", err)
	}

	for _, p := range pd {
		if p.IsDir() {
			continue
		}

		path := path.Join(pageDir, p.Name())
		rawData, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read file %s: %s", path, err)
		}
		data := strings.Split(string(rawData), "+++")
		name := strings.TrimSuffix(p.Name(), ".md")
		pages[name] = page{content: data[2]}
	}

	// Incoming content
	plugins := make(map[string]page)

	pl, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		log.Fatalf("Failed to get list of plugins: %s", err)
	}

	for _, p := range pl {
		if !p.IsDir() {
			continue
		}
		if getBlacklist()[p.Name()] {
			continue
		}

		path := path.Join(pluginDir, p.Name(), "README.md")
		rawData, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read file %s: %s", path, err)
		}
		pageData, err := parseReadme(rawData)
		plugins[p.Name()] = pageData

	}

	// Compare current and incoming state
	// Return a negative list to delete
	// Return a positive list to create or soverwrite
	neg, pos := compare(plugins, pages)

	if neg == nil && pos == nil {
		log.Printf("Nothing to update.")
	}

	for _, p := range neg {
		path := path.Join(pageDir, p)
		err := os.Remove(path + ".md")
		if err != nil {
			log.Fatalf("Error deleting file file %s: %s", path, err)
		}
	}

	for _, p := range pos {
		h := plugins[p].header
		h.Date = time.Now().UTC()
		h.Tags = []string{"plugin", h.Title}
		h.Categories = []string{"plugin"}
		h.Weight = 0 //TODO: predefined !sorting! see python code

		var buf bytes.Buffer
		err := toml.NewEncoder(&buf).Encode(h)
		if err != nil {
			log.Fatalf("Error creating TOML header data: %s", err)
		}

		headerData := buf.String()
		contentData := plugins[p].content
		path := path.Join(pageDir, p)
		fileContent := strings.Join([]string{"", headerData, contentData}, "+++\n")
		err = ioutil.WriteFile(path+".md", []byte(fileContent), 0644)
		if err != nil {
			log.Fatalf("Error writing file: %s", err)
		}
	}
}

func parseReadme(data []byte) (page, error) {
	var h meta

	parser := parser.New()
	doc := parser.Parse(data)

	// Get header data from README data
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}

		switch node.(type) {
		case *ast.Heading:
			if getString(ast.GetFirstChild(node)) == "Name" {
				c := ast.GetNextNode(node).GetChildren()
				h.Title = getString(ast.GetFirstChild(c[1]))
				description := strings.Split(getString(c[2]), " - ")[1]
				h.Description = fmt.Sprintf("*%s* %s", h.Title, description)

				// Remove header only nodes from tree
				ast.RemoveFromTree(ast.GetPrevNode(node))
				ast.RemoveFromTree(ast.GetNextNode(node))
				ast.RemoveFromTree(node)
				return ast.Terminate
			}
		}
		return ast.GoToNext
	})
	p := page{header: h, content: string(markdown.Render(doc, md.NewRenderer()))}
	return p, nil
}

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

func getString(node ast.Node) string {
	if n := node.AsContainer(); n != nil {
		return string(n.AsContainer().Literal)
	}

	if n := node.AsLeaf(); n != nil {
		return string(n.AsLeaf().Literal)
	}
	return ""
}

func getContent(node ast.Node) string {
	if n := node.AsContainer(); n != nil {
		return string(n.AsContainer().Content)
	}

	if n := node.AsLeaf(); n != nil {
		return string(n.AsLeaf().Content)
	}
	return ""
}

func getBlacklist() map[string]bool {
	var bl = map[string]bool{
		"test":       true,
		"pkg":        true,
		"deprecated": true,
	}
	return bl
}
