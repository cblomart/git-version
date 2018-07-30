package main

import (
	"log"

	git "gopkg.in/src-d/go-git.v4"
)

func main() {
	path := "./"
	// open git repo
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Printf("Could not open git repository '%s'\n", path)
	}
	// get logs
	l, err := r.Log(&git.LogOptions{})
	if err != nil {
		log.Printf("Could not get logs from repository '%s'\n", path)
	}
	// get last commit
	c, err := l.Next()
	if err != nil {
		log.Printf("Could not get last commit")
	}
	// print info on the commit
	log.Printf("Commit: %s\n", c.Hash.String())
}
