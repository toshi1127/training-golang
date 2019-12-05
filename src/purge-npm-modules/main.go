package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WalkTargetDir(path string, name string, isPrintOnly bool) error {
	// ioutil.ReadDir
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed: read directory %v: %v", path, err)
	}
	for _, file := range contents {
		absolutePath := filepath.Join(path, file.Name())

		if file.IsDir() && file.Name() == name {
			if isPrintOnly {
				fmt.Fprintln(os.Stdout, absolutePath)
			} else {
				if err := os.RemoveAll(absolutePath); err != nil {
					return fmt.Errorf("failed: remove directory %s: %v", absolutePath, err)
				}
			}

			fmt.Fprintln(os.Stdout, absolutePath)
			continue
		}

		if file.IsDir() {
			if err := WalkTargetDir(absolutePath, name, isPrintOnly); err != nil {
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
	isPrintOnly := flag.Bool("dir", false, "output found directories only - do not remove")
	flag.Parse()

	path := "."
	args := flag.Args()
	if len(args) > 0 && args[0] != "" {
		path = args[0]
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %s: %v\n", path, err)
		os.Exit(errorParseExitCode)
	}

	err = WalkTargetDir(absolutePath, "node_modules", *isPrintOnly)
	err = WalkTargetDir(absolutePath, "bundle", *isPrintOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v\n", err)
		os.Exit(errorExitCode)
	}
}