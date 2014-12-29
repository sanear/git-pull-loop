package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"os/exec"
	"sync"
)

func main() {
	wd, _ := os.Getwd()
	dirs, _ := ioutil.ReadDir(wd)
	didAnything := false
	wg := new(sync.WaitGroup)
	wg.Add(len(dirs))
	for _, dir := range dirs {
		if dir.IsDir() {
			if isRepo(dir.Name()) {
				func(dirName string, wg *sync.WaitGroup) {
					gitPull(dirName)
					defer func() {
						didAnything = true
						wg.Done()
					}()
				}(dir.Name(), wg)
			} else {
				wg.Done()
			}
		} else {
			wg.Done()
		}
	}	
	wg.Wait()
	if !didAnything {
		fmt.Println("No git repositories found")
	}
}

// Perform a git pull in the current working directory
// (assumes that git is on your $PATH)
func gitPull(dir string) {
	cmd := exec.Command("git", "-C", dir, "pull")
	fmt.Println("Executing ", cmd.Args, "...")
	err := cmd.Run()
	if nil != err {
		fmt.Println("Failed to pull from ", dir)
		return
	}
	fmt.Println("Finished pull from ", dir)
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
// func traverseAndDo(dir string) []string {
// 	files, _ := ioutil.ReadDir(dir)
// 	var subdirs []string
// 	for _, file := range files {
// 		if file.IsDir && !isRepo(file) {
// 			append(subdirs, file.Name())
// 		}
// 	}

// 	for _, subdir := range subdirs {
// 		traverse(subdir)
// 	}
// }
