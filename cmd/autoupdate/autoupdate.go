package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	workdir, _ := os.Getwd()
	filepath.WalkDir(".", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			if d.Name() == "vendor" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			os.Chdir(path)
			log.Println("entering", path)
			exec.Command("go", "get", "-u").Run()
			log.Println("leaving", path)
			os.Chdir(workdir)
		}
		return nil
	})
	os.Chdir(workdir)
}
