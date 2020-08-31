package main

import (
	"os"
	pathUtils "path"
	"strings"
)

/**
scan a folder to find any sub-folders
and git repositories.
*/
func Scan(f string) []string {
	return scanFolder(make([]string, 0), f)
}


/**
scan a given folder.
Retrieves all the sub-folders which are git repositories.
*/
func scanFolder(paths []string, curr string) []string {
	f, _ := os.Open(curr)

	fps, _ := f.Readdir(-1)
	defer f.Close()

	var path string
	for _, fp := range fps {
		//for each directory, find out if is a git repo
		if fp.IsDir(){
			path = pathUtils.Join(curr, fp.Name())
			if fp.Name() == ".git" {
				path = strings.TrimSuffix(path, "\\.git")
				paths = append(paths, path)
				continue
			}

			//ignore some common paths which are of no use
			if fp.Name() == "vendor" || fp.Name() == "node_modules" {
				continue
			}

			paths = scanFolder(paths, path)
		}
	}

	return paths
}


// User represents a user account.
type User struct {
	// Uid is the user ID.
	Uid string
	// Gid is the primary group ID.
	Gid string
	// Username is the login name.
	Username string
	// Name is the user's real or display name.
	Name string
	// HomeDir is the path to the user's home directory (if they have one).
	HomeDir string
}