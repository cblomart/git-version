package main

//go:generate git-version

import (
	"log"
	"os"
	"strings"
	"text/template"

	git "gopkg.in/src-d/go-git.v4"
)

const (
	versionFilename string = "./version.go"
	versionTemplate string = `
package main

const(
	gitCommit = "{{or .Commit "-"}}"
	gitShortCommit = "{{or .ShortCommit "-"}}"
	gitTag= "{{or .Tag "-"}}"
	gitStatus = "{{or .Status "-"}}" 
)
`
)

//GitInfo has the informations reported by git-verion
type GitInfo struct {
	Commit      string
	ShortCommit string
	Tag         string
	Status      string
}

func main() {
	// check current repository
	path := "./"
	// create the info object
	infos := GitInfo{}
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
	infos.Commit = commit.Hash.String()
	infos.ShortCommit = infos.Commit[:7]
	// print info on the commit
	log.Printf("Commit: %s\n", infos.Commit)
	log.Printf("Short Commit: %s\n", infos.ShortCommit)
	// get tags
	tag, err := repo.TagObject(commit.Hash)
	if err == nil {
		infos.Tag = strings.Replace(tag.Name, "refs/tags/", "", 1)
	}
	log.Printf("Tag: %s\n", infos.Tag)
	// get status
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
	infos.Status = "clean"
	if !status.IsClean() {
		infos.Status = "dirty"
	}
	log.Printf("Status: %s\n", infos.Status)
	// generate version.go file
	t := template.New("version.go")
	_, err = t.Parse(versionTemplate)
	if err != nil {
		log.Printf("Could not parse template: %s", err)
		return
	}
	// create file
	f, err := os.Create(versionFilename)
	defer f.Close()
	if err != nil {
		log.Printf("Could not create file %s: %s\n", versionFilename, err)
		return
	}
	t.Execute(f, infos)
}
