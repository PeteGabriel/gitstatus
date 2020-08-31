package main

import (
	"flag"
)

func main() {
	var folder string
	var email string

	flag.StringVar(&folder, "src", "", "source folder to scan for git repositories")
	flag.StringVar(&email, "email", "", "the email to scan")

	flag.Parse()

	var repos []string
	if folder != "" {
		repos = Scan(folder)
		return
	}

	Stats(email, repos)
}
