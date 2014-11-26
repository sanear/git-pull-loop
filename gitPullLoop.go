package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"os/exec"
)

func main() {
	wd, _ := os.Getwd()
	dirs, _ := ioutil.ReadDir(wd)
	didAnything := false
	for _, dir := range dirs {
		if dir.IsDir() {
			os.Chdir(dir.Name())
			if isRepo("./") {
				fmt.Println(dir.Name())
				go gitPull() 	// Run the long-taking part on its own thread
				didAnything = true
			}
			os.Chdir("../")
		}
	}
	if !didAnything {
		fmt.Println("No git repositories found")
	}
}

// Perform a git pull in the current working directory
// (assumes that git is on your $PATH)
func gitPull() {
	cmd := exec.Command("git", "pull")
	err := cmd.Run()
	if nil != err {
		fmt.Println(err)
		panic(err)
	}	
}

// Given a relative path as a string, determine if the directory contains
// a .git directory
func isRepo(dir string) bool {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if  file.Name() == ".git" {
			return true
		}
	}
	return false
}

// Starting from a string directory, recurse to find all directories
// containing .git directories (i.e., find all git repos)
func traverseAndDo(dir string) []string {
	files, _ := ioutil.ReadDir(dir)
	var subdirs []string
	for _, file := range files {
		if file.IsDir && !isRepo(file) {
			append(subdirs, file.Name())
		}
	}

	for _, subdir := range subdirs {
		traverse(subdir)
	}
}
