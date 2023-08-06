package main

//go:generate git-version

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	originName      string = "origin"
	versionFilename string = "./version.go"
	versionTemplate string = `//lint:file-ignore U1000 Ignore all unused code as it is generated
package {{.Package}}

const (
	{{if .Expose}}G{{else}}g{{end}}itCommit      = "{{if .Commit}}{{.Commit}}{{else}}NA{{end}}"
	{{if .Expose}}G{{else}}g{{end}}itShortCommit = "{{if .ShortCommit}}{{.ShortCommit}}{{else}}NA{{end}}"
	{{if .Expose}}G{{else}}g{{end}}itTag         = "{{if .Tag}}{{.Tag}}{{else}}NA{{end}}"
	{{if .Expose}}G{{else}}g{{end}}itBranch      = "{{if .Branch}}{{.Branch}}{{else}}NA{{end}}"
	{{if .Expose}}G{{else}}g{{end}}itStatus      = "{{if .Status}}{{.Status}}{{else}}NA{{end}}"
)
`
)

// GitInfo has the informations reported by git-verion
type GitInfo struct {
	Expose      bool
	Package     string
	Commit      string
	ShortCommit string
	Tag         string
	Branch      string
	Status      string
}

func main() {
	// argument parsing
	var gitPath string
	var outputFilePath string
	var goPackage string
	var expose bool

	// declare arguments
	flag.StringVar(&gitPath, "g", "./", "Git path. Defaults to the current directory")
	flag.StringVar(&outputFilePath, "o", "./version.go", "Output file path. Directories will be created if they don't exist. Defauls to './version.go'")
	flag.StringVar(&goPackage, "p", "main", "Package for the output. Defauls to 'main'")
	flag.BoolVar(&expose, "e", false, "Expose version globaly. Disabled by default")

	// parse the arguments
	flag.Parse()

	// create the info object
	infos := GitInfo{Package: goPackage, Expose: expose}
	// open git repo
	repo, err := git.PlainOpen(gitPath)
	if err != nil {
		log.Printf("Could not open git repository '%s'\n", gitPath)
		return
	}
	// get logs from local
	logs, err := repo.Log(&git.LogOptions{})
	if err != nil {
		log.Printf("Could not get logs from repository '%s'\n", gitPath)
		return
	}
	// get last commit from local
	commit, err := logs.Next()
	if err != nil {
		log.Printf("Could not get last commit from repository '%s'\n", gitPath)
		return
	}
	infos.Commit = commit.Hash.String()
	infos.ShortCommit = infos.Commit[:7]
	// print info on the commit
	log.Printf("Commit: %s\n", infos.Commit)
	log.Printf("Short Commit: %s\n", infos.ShortCommit)
	// get the branch
	head, err := repo.Head()
	if err != nil {
		log.Println("Could not get head")
		return
	}
	infos.Branch = strings.Replace(string(head.Name()), "refs/heads/", "", 1)
	log.Printf("Branch: %s\n", infos.Branch)
	// fetch tags
	err = repo.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/tags/*:refs/tags/*"},
	})
	if err != nil {
		log.Printf("could not fetch tags")
	}
	// get tags
	iter, err := repo.Tags()
	if err != nil {
		log.Printf("Could not get tags from repository '%s'\n", gitPath)
	}
	iter.ForEach(func(r *plumbing.Reference) error {
		if r.Hash().String() == infos.Commit {
			infos.Tag = strings.Replace(string(r.Name()), "refs/tags/", "", 1)
			return nil
		}
		return nil
	})
	log.Printf("Tag: %s\n", infos.Tag)
	// get status
	worktree, err := repo.Worktree()
	if err != nil {
		log.Printf("Couldn't get the worktree from repository '%s'", gitPath)
		return
	}
	status, err := worktree.Status()
	if err != nil {
		log.Printf("Couldn't get the worktree from repository '%s'", gitPath)
		return
	}
	infos.Status = "clean"
	if !status.IsClean() {
		infos.Status = "dirty"
	}
	log.Printf("Status: %s\n", infos.Status)
	// generate version.go file
	t := template.New(filepath.Base(outputFilePath))
	_, err = t.Parse(versionTemplate)
	if err != nil {
		log.Printf("Could not parse template: %s", err)
		return
	}
	// create outputdir if needed
	err = os.MkdirAll(filepath.Dir(outputFilePath), 0777)
	if err != nil && !os.IsExist(err) {
		log.Printf("Could not create output directory: %s", err)
		return
	}
	// create file
	f, err := os.Create(outputFilePath)
	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Could not close file %s: %s\n", outputFilePath, err)
			return
		}
	}()
	if err != nil {
		log.Printf("Could not create file %s: %s\n", versionFilename, err)
		return
	}
	t.Execute(f, infos)
}
