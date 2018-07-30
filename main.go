package main

import (
	"log"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
)

func main() {
	path := "./"
	// open git repo
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Printf("Could not open git repository '%s'\n", path)
		return
	}
	// get logs
	logs, err := repo.Log(&git.LogOptions{})
	if err != nil {
		log.Printf("Could not get logs from repository '%s'\n", path)
		return
	}
	// get last commit
	commit, err := logs.Next()
	if err != nil {
		log.Printf("Could not get last commit from repository '%s'\n", path)
	}
	commitname := commit.Hash.String()
	// print info on the commit
	log.Printf("Commit: %s\n", commitname)
	// get tags
	tagname := ""
	tag, err := repo.TagObject(commit.Hash)
	if err == nil {
		tagname = strings.Replace(tag.Name, "refs/tags/", "", 1)
	}
	log.Printf("Tag: %s\n", tagname)
	worktree, err := repo.Worktree()
	if err != nil {
		log.Printf("Couldn't get the worktree from repository '%s'", path)
		return
	}
	status, err := worktree.Status()
	if err != nil {
		log.Printf("Couldn't get the worktree from repository '%s'", path)
		return
	}
	gitstatus := "clean"
	if !status.IsClean() {
		gitstatus = "dirty"
		log.Println(status.String())
	}
	log.Printf("Status: %s\n", gitstatus)
}
