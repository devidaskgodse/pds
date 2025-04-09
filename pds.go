package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
)

type PathType string

const (
	File PathType = "file"
	Dir  PathType = "dir"
)

type PathEntry struct {
	Path string
	Type PathType
}

// Flatten the TOML structure and build a unique key map
func buildPathDict(tree *toml.Tree, base string, keyParts []string, paths map[string]PathEntry) {
	for _, key := range tree.Keys() {
		val := tree.Get(key)
		currentPath := filepath.Join(base, key)
		keyPath := append(keyParts, key)
		mapKey := pathToKey(keyPath)

		switch v := val.(type) {
		case *toml.Tree:
			// Directory
			paths[mapKey] = PathEntry{Path: currentPath, Type: Dir}
			buildPathDict(v, currentPath, keyPath, paths)
		case []interface{}:
			if key == "files" {
				for _, file := range v {
					if fileStr, ok := file.(string); ok {
						fullPath := filepath.Join(base, fileStr)
						fileKey := pathToKey(append(keyParts, strings.TrimSuffix(fileStr, filepath.Ext(fileStr))))
						paths[fileKey] = PathEntry{Path: fullPath, Type: File}
					}
				}
			}
		}
	}
}

func pathToKey(parts []string) string {
	joined := strings.Join(parts, "_")
	return strings.ReplaceAll(joined, ".", "")
}

func createStructure(paths map[string]PathEntry) error {
	for _, entry := range paths {
		switch entry.Type {
		case Dir:
			err := os.MkdirAll(entry.Path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory '%s': %w", entry.Path, err)
			}
			fmt.Printf("Created directory: %s\n", entry.Path)
		case File:
			err := os.MkdirAll(filepath.Dir(entry.Path), os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating parent directory for '%s': %w", entry.Path, err)
			}
			f, err := os.Create(entry.Path)
			if err != nil {
				return fmt.Errorf("error creating file '%s': %w", entry.Path, err)
			}
			f.Close()
			fmt.Printf("Created file: %s\n", entry.Path)
		}
	}
	return nil
}

func main() {
	// CLI flags
	tomlPath := flag.String("config", "structure.toml", "Path to the TOML structure file")
	baseDir := flag.String("out", "new_project", "Base output directory where structure will be created")
	flag.Parse()

	tree, err := toml.LoadFile(*tomlPath)
	if err != nil {
		log.Fatalf("Failed to load TOML file %s: %v", *tomlPath, err)
	}

	paths := make(map[string]PathEntry)
	buildPathDict(tree, *baseDir, []string{}, paths)

	err = createStructure(paths)
	if err != nil {
		log.Fatalf("Error creating structure: %v", err)
	}

	fmt.Println("\nNamed Path Dictionary:")
	for key, entry := range paths {
		fmt.Printf("%-25s => %s [%s]\n", key, entry.Path, entry.Type)
	}
}

