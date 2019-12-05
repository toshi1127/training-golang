package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WalkTargetDir(path string, name string, printOnly bool) error {
	// ioutil.ReadDir
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read content of directory %v: %v", path, err)
	}
	for _, file := range contents {
		absolutePath := filepath.Join(path, file.Name())

		if file.IsDir() && file.Name() == name {
			if printOnly {
				fmt.Fprintln(os.Stdout, absolutePath)
			} else {
				if err := os.RemoveAll(absolutePath); err != nil {
					return fmt.Errorf("failed to remove directory %s: %v", absolutePath, err)
				}
			}

			fmt.Fprintln(os.Stdout, absolutePath)
			continue
		}

		if file.IsDir() {
			if err := WalkTargetDir(absolutePath, name, printOnly); err != nil {
				return err
			}
		}
	}
	return nil
}

var (
	errorExitCode      = 1
	errorParseExitCode = 2
)

func main() {
	dirName := flag.Bool("dir", false, "output found directories only - do not remove")
	flag.Parse()

	path := "."
	args := flag.Args()
	if len(args) > 0 && args[0] != "" {
		path = args[0]
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse given path %s: %v\n", path, err)
		os.Exit(errorParseExitCode)
	}

	err = WalkTargetDir(absolutePath, "node_modules", *dirName)
	err = WalkTargetDir(absolutePath, "bundle", *dirName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "purging failed with an error: %v\n", err)
		os.Exit(errorExitCode)
	}
}