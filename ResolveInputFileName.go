package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type ResolveInputFileName struct {
}

func (self *ResolveInputFileName) run(inputFileName string) (string, error) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	workingDir, _ := os.Getwd()
	if inputFileName == "" {
		return inputFileName, nil
	}

	if !filepath.IsAbs(inputFileName) {
		if strings.HasPrefix(inputFileName, "~/") {
			inputFileName = filepath.Join(homeDir, inputFileName[2:])
		} else {
			inputFileName = filepath.Join(workingDir, inputFileName)
		}
	}
	return inputFileName, nil
}
